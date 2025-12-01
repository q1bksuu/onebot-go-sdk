"""
Markdown 文档解析器
解析 OneBot 标准 Markdown 文档，提取 API 定义、事件定义等信息
"""

import re
from typing import List, Optional, Dict, Tuple
from pathlib import Path
from schema import Field, APIModel, APIDefinition, EventModel, FieldType, MessageTypeVariant, MessageSegment
from type_mapper import TypeMapper


class MarkdownParser:
    """Markdown 文档解析器"""

    def __init__(self):
        self.type_mapper = TypeMapper()

    def parse_api_file(self, file_path: str) -> List[APIDefinition]:
        """
        解析 API Markdown 文件，提取所有 API 定义

        Args:
            file_path: Markdown 文件路径

        Returns:
            API 定义列表
        """
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()

        return self.extract_apis_from_content(content)

    def extract_apis_from_content(self, content: str) -> List[APIDefinition]:
        """从 Markdown 内容中提取 API 定义"""
        apis = []

        # 正则：匹配 ## `api_name` 描述
        api_pattern = r"^##\s+`(\w+)`\s+(.+?)$"

        lines = content.split("\n")
        i = 0

        while i < len(lines):
            line = lines[i]
            match = re.match(api_pattern, line)

            if match:
                api_name = match.group(1)
                description = match.group(2).strip()

                # 当前位置是 API 的开始，下一行开始查找参数表
                i += 1

                # 提取参数表（从当前位置向后查找 ### 参数）
                params_model = self._extract_model_from_lines(lines, i, "参数", api_name + "_req")

                # 提取响应数据表（从当前位置向后查找 ### 响应数据）
                response_model = self._extract_model_from_lines(lines, i, "响应数据", api_name + "_resp")

                # 创建 API 定义
                api_def = APIDefinition(
                    api_name=api_name,
                    description=description,
                    request_model=params_model,
                    response_model=response_model,
                )
                apis.append(api_def)
            else:
                i += 1

        return apis

    def _extract_model_from_lines(self, lines: List[str], start_idx: int, section_name: str, model_name: str) -> APIModel:
        """
        从 Markdown 行列表中提取模型（参数或响应）

        Args:
            lines: Markdown 文件行列表
            start_idx: 开始位置
            section_name: 节点名称（如"参数"、"响应数据"）
            model_name: 模型名称

        Returns:
            APIModel 对象
        """
        fields = []

        # 查找对应的标题（### 参数 或 ### 响应数据）
        header_found = False
        i = start_idx

        while i < len(lines):
            if f"### {section_name}" in lines[i]:
                header_found = True
                i += 1
                break
            i += 1

        if not header_found:
            return APIModel(name=model_name, fields=[])

        # 跳过空行，找到表格开始
        while i < len(lines):
            line = lines[i].strip()
            if line.startswith("|"):
                # 找到表格开始
                break
            i += 1

        # 跳过表格头行 (| 字段名 | 数据类型 |...)
        if i < len(lines):
            i += 1

        # 跳过表格分隔行 (| --- | --- |...)
        if i < len(lines) and "-" in lines[i]:
            i += 1

        # 解析表格行
        while i < len(lines):
            line = lines[i].strip()

            # 表格结束条件
            if not line or not line.startswith("|"):
                break

            # 解析表格行
            field = self._parse_table_row(line, section_name)
            if field:
                fields.append(field)

            i += 1

        return APIModel(name=model_name, fields=fields)

    def _parse_table_row(self, row: str, section_name: str) -> Optional[Field]:
        """
        解析表格行，提取字段信息

        Args:
            row: 表格行字符串（|field|type|value|desc|）
            section_name: 节点名称

        Returns:
            Field 对象或 None
        """
        # 分割行
        parts = [p.strip() for p in row.split("|")]
        parts = [p for p in parts if p]  # 移除空部分

        # 期望至少 3 列
        if len(parts) < 3:
            return None

        field_name = parts[0]
        data_type = parts[1]
        description = parts[-1]  # 最后一列总是描述

        # 获取默认值或可能的值（取决于列数和上下文）
        default_value = None
        possible_values = []
        if len(parts) >= 4:
            middle_column = parts[2]
            # 如果中间列包含中文顿号、逗号或反引号包裹的值，视为可能的值
            # 单个反引号包裹的值也算作可能的值（如 `private`）
            if "、" in middle_column or "，" in middle_column or ("`" in middle_column and middle_column.count("`") >= 2):
                # 解析可能的值
                possible_values = self._parse_possible_values(middle_column)
            elif middle_column and middle_column != "-":
                # 否则视为默认值
                default_value = middle_column

        # 跳过无效的表头行
        if field_name.startswith("-") or "字段名" in field_name or "名" in field_name:
            return None

        # 清理字段名：移除 Markdown 代码格式的反引号 `field_name`
        field_name = field_name.strip("`")

        # 处理含有 "或" 的字段名，只取第一个
        # 例如: `anonymous_flag` 或 `flag` -> anonymous_flag
        if " 或 " in field_name:
            field_name = field_name.split(" 或 ")[0].strip("`")

        # 跳过特殊字符或无效的字段名
        if not field_name or field_name.startswith("……") or field_name == "……" or not field_name.replace("_", "").isalnum():
            return None

        # 清理数据类型和默认值中的反引号
        data_type = data_type.strip("`")
        if default_value:
            default_value = default_value.strip("`")

        # 解析数据类型
        field_type, go_type = self.type_mapper.parse_type(data_type)

        # 确定是否必需
        required = self.type_mapper.is_required_field(default_value)
        use_omitempty = not required or (default_value is not None and default_value.strip() != "")

        # 特殊处理 message 类型
        message_variants = []
        if field_type == FieldType.MESSAGE:
            message_variants = self.type_mapper.get_message_variants()

        # 处理 Go 字段名
        go_field_name = self.type_mapper.snake_to_pascal(field_name)

        return Field(
            name=field_name,
            go_name=go_field_name,
            data_type=data_type,
            field_type=field_type,
            go_type=go_type,
            description=description,
            required=required,
            default_value=default_value,
            is_optional=use_omitempty,
            message_variants=message_variants,
            possible_values=possible_values,
        )

    def _parse_possible_values(self, value_str: str) -> List[str]:
        """
        解析可能的值字符串

        Args:
            value_str: 值字符串，例如 "`friend`、`group`、`other`" 或 "friend, group, other"

        Returns:
            值列表
        """
        if not value_str or value_str == "-":
            return []

        # 替换中文标点为英文标点
        value_str = value_str.replace("、", ",").replace("，", ",")

        # 分割并清理
        values = [v.strip().strip("`").strip() for v in value_str.split(",")]
        values = [v for v in values if v and v != "-"]

        return values

    def _skip_to_next_table(self, lines: List[str], start_idx: int) -> int:
        """跳转到下一个表格的开始位置"""
        i = start_idx
        while i < len(lines):
            if "|" in lines[i] and "---" not in lines[i]:
                return i
            i += 1
        return i

    def parse_event_file(self, file_path: str) -> List[EventModel]:
        """
        解析事件 Markdown 文件

        Args:
            file_path: Markdown 文件路径

        Returns:
            事件定义列表
        """
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()

        return self.extract_events_from_content(content)

    def extract_events_from_content(self, content: str) -> List[EventModel]:
        """从 Markdown 内容中提取事件定义"""
        events = []

        # 正则：匹配 ## 事件名称
        event_pattern = r"^##\s+(.+?)$"

        lines = content.split("\n")
        i = 0

        while i < len(lines):
            line = lines[i]
            match = re.match(event_pattern, line)

            if match:
                event_name = match.group(1).strip()

                # 跳过目录标题（如"# 消息事件"、"# 通知事件"等）
                if event_name.startswith("#") or "目录" in event_name:
                    i += 1
                    continue

                # 当前位置是事件的开始，下一行开始查找事件数据表
                i += 1

                # 提取事件数据模型
                event_model_dict = self._extract_event_model_from_lines(lines, i, event_name)

                if event_model_dict:
                    # 将字典转换为 EventModel 对象
                    event_model = EventModel(
                        name=event_model_dict["name"],
                        fields=event_model_dict["fields"],
                        event_type=event_model_dict["event_type"],
                        sub_type=event_model_dict["sub_type"],
                        description=event_model_dict["description"],
                    )
                    events.append(event_model)
            else:
                i += 1

        return events

    def _extract_event_model_from_lines(self, lines: List[str], start_idx: int, event_name: str) -> Optional[Dict]:
        """
        从 Markdown 行列表中提取事件模型

        Args:
            lines: Markdown 文件行列表
            start_idx: 开始位置
            event_name: 事件名称

        Returns:
            事件模型字典或 None
        """
        # 查找 "### 事件数据" 或 "### 上报数据" 标题，或直接查找表格
        header_found = False
        i = start_idx

        while i < len(lines):
            line = lines[i].strip()
            if "### 事件数据" in line or "### 上报数据" in line:
                header_found = True
                i += 1
                break
            # 如果直接遇到表格，也认为找到了事件数据
            if line.startswith("|") and "字段名" in line:
                header_found = True
                break
            # 如果遇到下一个事件定义，说明当前事件没有数据表
            if line.startswith("## "):
                return None
            i += 1

        if not header_found:
            return None

        # 跳过空行，找到表格开始
        while i < len(lines):
            line = lines[i].strip()
            if line.startswith("|"):
                break
            if line.startswith("## "):
                return None
            i += 1

        # 跳过表格头行
        if i < len(lines):
            i += 1

        # 跳过表格分隔行
        if i < len(lines) and "-" in lines[i]:
            i += 1

        # 解析事件主表格
        fields = []
        nested_objects = {}  # 存储嵌套对象的字段信息

        while i < len(lines):
            line = lines[i].strip()

            # 表格结束条件
            if not line or not line.startswith("|"):
                break

            # 解析表格行
            field = self._parse_table_row(line, "事件数据")
            if field:
                fields.append(field)

                # 如果字段类型是 object，查找其嵌套字段定义
                if field.field_type == FieldType.OBJECT:
                    nested_fields = self._extract_nested_object_fields(lines, i + 1, field.name)
                    if nested_fields:
                        nested_objects[field.name] = nested_fields

            i += 1

        # 构建事件模型
        # 从字段中提取事件类型信息
        post_type = ""
        event_type = ""
        sub_type = ""

        for field in fields:
            if field.name == "post_type" and field.possible_values:
                post_type = field.possible_values[0] if len(field.possible_values) == 1 else ""
            elif field.name in ["message_type", "notice_type", "request_type", "meta_event_type"]:
                if field.possible_values and len(field.possible_values) > 0:
                    event_type = field.possible_values[0]
            elif field.name == "sub_type" and field.possible_values:
                sub_type = ",".join(field.possible_values)

        # 将嵌套对象的字段信息附加到对应的字段上
        for field in fields:
            if field.name in nested_objects:
                field.nested_fields = nested_objects[field.name]

        return {
            "name": event_name,
            "fields": fields,
            "event_type": event_type,
            "sub_type": sub_type,
            "description": event_name,
        }

    def _extract_nested_object_fields(self, lines: List[str], start_idx: int, object_name: str) -> List[Field]:
        """
        提取嵌套对象的字段定义

        Args:
            lines: Markdown 文件行列表
            start_idx: 开始位置
            object_name: 对象名称

        Returns:
            字段列表
        """
        i = start_idx

        # 查找类似 "其中 `object_name` 字段的内容如下：" 的说明
        while i < len(lines):
            line = lines[i].strip()
            if f"`{object_name}`" in line and ("字段的内容如下" in line or "内容如下" in line):
                i += 1
                break
            # 遇到下一个二级或三级标题，停止查找
            if line.startswith("##"):
                return []
            i += 1
            # 防止无限循环，最多向后查找30行
            if i - start_idx > 30:
                return []

        # 跳过空行，找到表格开始
        while i < len(lines):
            line = lines[i].strip()
            if line.startswith("|"):
                break
            if line.startswith("##"):
                return []
            i += 1

        # 跳过表格头行
        if i < len(lines):
            i += 1

        # 跳过表格分隔行
        if i < len(lines) and "-" in lines[i]:
            i += 1

        # 解析表格行
        fields = []
        while i < len(lines):
            line = lines[i].strip()

            if not line or not line.startswith("|"):
                break

            field = self._parse_table_row(line, "嵌套对象")
            if field:
                fields.append(field)

            i += 1

        return fields

    def parse_message_segment_file(self, file_path: str) -> List[MessageSegment]:
        """
        解析消息段 Markdown 文件

        Args:
            file_path: Markdown 文件路径

        Returns:
            消息段定义列表
        """
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()

        return self.extract_message_segments_from_content(content)

    def extract_message_segments_from_content(self, content: str) -> List[MessageSegment]:
        """从 Markdown 内容中提取消息段定义"""
        from schema import MessageSegment

        segments = []
        lines = content.split("\n")
        i = 0

        while i < len(lines):
            line = lines[i]

            # 匹配二级标题,如 "## 纯文本", "## QQ 表情" 等
            if line.startswith("## "):
                # 保存原始标题(用于检查 Badge 标记)
                original_title = line[3:].strip()

                # 提取标题内容,移除 Badge 等标记
                segment_title = re.sub(r'<Badge[^>]*>', '', original_title).strip()

                # 跳过目录级别的标题
                if not segment_title or segment_title.startswith("#"):
                    i += 1
                    continue

                description = segment_title
                segment_start = i  # 记录段落开始位置
                i += 1

                # 查找当前段落的结束位置(下一个 ## 标题或文件末尾)
                segment_end = len(lines)
                for j in range(i, len(lines)):
                    if lines[j].startswith("## "):
                        segment_end = j
                        break

                # 在当前段落中查找 JSON 示例以获取消息段类型
                segment_type = None
                for j in range(segment_start, segment_end):
                    line_content = lines[j].strip()
                    # 尝试从 JSON 示例中提取 type 字段
                    type_match = re.search(r'"type":\s*"(\w+)"', line_content)
                    if type_match:
                        segment_type = type_match.group(1)
                        break

                # 如果没有找到类型,跳过这个段落
                if not segment_type:
                    i = segment_end
                    continue

                # 在当前段落中查找参数表
                fields = []
                can_send = False
                can_receive = False
                table_found = False

                j = segment_start
                while j < segment_end:
                    line_content = lines[j].strip()

                    # 找到表格头(包含"参数名"、"收"、"发"等)
                    if "|" in line_content and ("参数名" in line_content or "字段名" in line_content):
                        table_found = True
                        j += 1

                        # 跳过分隔行
                        if j < len(lines) and "-" in lines[j]:
                            j += 1

                        # 解析表格数据行
                        while j < segment_end:
                            data_line = lines[j].strip()

                            # 表格结束
                            if not data_line or not data_line.startswith("|"):
                                break

                            # 解析表格行,提取字段信息
                            field = self._parse_message_segment_table_row(data_line)
                            if field:
                                fields.append(field)

                                # 检查是否可发送/接收
                                parts = [p.strip() for p in data_line.split("|")]
                                parts = [p for p in parts if p]
                                if len(parts) >= 3:
                                    if len(parts) >= 2 and "✓" in parts[1]:
                                        can_receive = True
                                    if len(parts) >= 3 and "✓" in parts[2]:
                                        can_send = True

                            j += 1
                        break

                    j += 1

                # 检查原始标题中的 Badge 标记来确定发送/接收能力
                if 'text="发"' in original_title or 'text="发' in original_title:
                    can_send = True
                    can_receive = False
                elif 'text="收"' in original_title or 'text="收' in original_title:
                    can_send = False
                    can_receive = True
                else:
                    # 如果没有明确标记,根据表格内容确定
                    if not can_send and not can_receive and table_found:
                        can_send = True
                        can_receive = True
                    elif not can_send and not can_receive:
                        # 如果既没有表格也没有标记,默认都支持
                        can_send = True
                        can_receive = True

                # 创建消息段对象
                segment = MessageSegment(
                    segment_type=segment_type,
                    fields=fields,
                    description=description,
                    can_send=can_send,
                    can_receive=can_receive,
                )
                segments.append(segment)

                # 跳转到下一个段落
                i = segment_end
                continue

            i += 1

        return segments

    def _parse_message_segment_table_row(self, row: str) -> Optional[Field]:
        """
        解析消息段参数表格行

        Args:
            row: 表格行字符串(如 | file | ✓ | ✓ | - | 图片文件名 |)

        Returns:
            Field 对象或 None
        """
        # 分割行
        parts = [p.strip() for p in row.split("|")]
        parts = [p for p in parts if p]  # 移除空部分

        # 期望至少 4 列: 参数名 | 收 | 发 | 说明
        # 或 5 列: 参数名 | 收 | 发 | 可能的值 | 说明
        if len(parts) < 4:
            return None

        field_name = parts[0].strip("`")

        # 跳过表头行
        if field_name.startswith("-") or "参数名" in field_name or "字段名" in field_name:
            return None

        # 根据列数确定各字段位置
        if len(parts) == 4:
            # | 参数名 | 收 | 发 | 说明 |
            description = parts[3]
            possible_values = []
        else:
            # | 参数名 | 收 | 发 | 可能的值 | 说明 |
            possible_values_str = parts[3]
            description = parts[4] if len(parts) > 4 else ""

            # 解析可能的值
            possible_values = self._parse_possible_values(possible_values_str)

        # 对于消息段参数,大多数类型都是 string
        # 特殊处理一些明确的类型
        data_type = "string"
        if field_name in ["id", "type"] and possible_values:
            # 如果有可能的值,保持为 string
            data_type = "string"

        field_type, go_type = self.type_mapper.parse_type(data_type)

        # 消息段的字段通常都是可选的(发送时)
        required = False
        use_omitempty = True

        # 处理 Go 字段名
        go_field_name = self.type_mapper.snake_to_pascal(field_name)

        return Field(
            name=field_name,
            go_name=go_field_name,
            data_type=data_type,
            field_type=field_type,
            go_type=go_type,
            description=description,
            required=required,
            default_value=None,
            is_optional=use_omitempty,
            possible_values=possible_values,
        )
