"""
Markdown 文档解析器
解析 OneBot 标准 Markdown 文档，提取 API 定义、事件定义等信息
"""

import re
from typing import List, Optional, Dict, Tuple
from pathlib import Path
from schema import Field, APIModel, APIDefinition, FieldType, MessageTypeVariant
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

        # 获取默认值（取决于列数）
        default_value = None
        if len(parts) >= 4:
            default_value = parts[2]

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
        )

    def _skip_to_next_table(self, lines: List[str], start_idx: int) -> int:
        """跳转到下一个表格的开始位置"""
        i = start_idx
        while i < len(lines):
            if "|" in lines[i] and "---" not in lines[i]:
                return i
            i += 1
        return i

    def parse_event_file(self, file_path: str) -> List[Dict]:
        """
        解析事件 Markdown 文件

        Args:
            file_path: Markdown 文件路径

        Returns:
            事件定义列表
        """
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()

        # TODO: 实现事件解析逻辑
        return []

    def parse_message_segment_file(self, file_path: str) -> List[Dict]:
        """
        解析消息段 Markdown 文件

        Args:
            file_path: Markdown 文件路径

        Returns:
            消息段定义列表
        """
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()

        # TODO: 实现消息段解析逻辑
        return []
