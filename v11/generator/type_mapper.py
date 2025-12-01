"""
类型系统映射
将 Markdown 中的数据类型字符串映射到 Go 类型
"""

import re
from typing import Tuple, Optional, List
from schema import FieldType, MessageTypeVariant


class TypeMapper:
    """类型映射器"""

    # 基础类型映射表
    TYPE_MAPPING = {
        # 数字类型
        r"number\s*\(\s*int64\s*\)": (FieldType.INT64, "int64"),
        r"number\s*\(\s*int32\s*\)": (FieldType.INT32, "int32"),
        r"number\s*\(\s*uint64\s*\)": (FieldType.UINT64, "uint64"),
        r"number\s*\(\s*uint32\s*\)": (FieldType.UINT32, "uint32"),
        r"number\s*\(\s*int\s*\)": (FieldType.INT, "int"),
        r"number\s*\(\s*uint\s*\)": (FieldType.UINT, "uint"),
        r"number": (FieldType.INT64, "int64"),  # 默认 number -> int64
        r"string": (FieldType.STRING, "string"),
        r"boolean": (FieldType.BOOL, "bool"),
        r"bool": (FieldType.BOOL, "bool"),

        # 复合类型
        r"object": (FieldType.OBJECT, "map[string]interface{}"),
        r"array": (FieldType.ARRAY, "[]interface{}"),

        # OneBot 特殊类型
        r"message": (FieldType.MESSAGE, "MessageValue"),  # 特殊处理
    }

    @staticmethod
    def parse_type(data_type_str: str) -> Tuple[FieldType, str]:
        """
        解析数据类型字符串，返回 (FieldType, Go类型)

        Args:
            data_type_str: 数据类型字符串，如 "number (int64)", "string", "object" 等

        Returns:
            (FieldType, Go类型字符串)
        """
        # 归一化字符串
        data_type_str = data_type_str.strip().lower()

        # 尝试匹配类型映射表
        for pattern, (field_type, go_type) in TypeMapper.TYPE_MAPPING.items():
            if re.match(f"^{pattern}$", data_type_str):
                return field_type, go_type

        # 如果没有匹配，返回 UNKNOWN
        return FieldType.UNKNOWN, "interface{}"

    @staticmethod
    def snake_to_camel(snake_str: str) -> str:
        """将 snake_case 转换为 CamelCase（Go 字段命名规则）"""
        components = snake_str.split("_")
        # 首个单词小写，后续单词首字母大写
        return components[0] + "".join(x.title() for x in components[1:])

    @staticmethod
    def snake_to_pascal(snake_str: str) -> str:
        """将 snake_case 转换为 PascalCase（Go 导出字段命名规则）"""
        components = snake_str.split("_")
        return "".join(x.title() for x in components)

    @staticmethod
    def is_message_type(data_type_str: str) -> bool:
        """判断是否是 message 类型"""
        return TypeMapper.parse_type(data_type_str)[0] == FieldType.MESSAGE

    @staticmethod
    def get_message_variants() -> List[MessageTypeVariant]:
        """
        获取 message 类型的可能变体
        OneBot 中 message 字段可以是：
        1. 字符串 (CQ 码格式)
        2. 消息段数组
        """
        return [MessageTypeVariant.STRING, MessageTypeVariant.ARRAY]

    @staticmethod
    def is_required_field(default_value: Optional[str]) -> bool:
        """
        判断字段是否是必需的
        如果有默认值，则字段不是必需的
        """
        return default_value is None or default_value.strip() == ""

    @staticmethod
    def determine_go_type_with_omitempty(
        field_type: FieldType, go_type: str, required: bool, default_value: Optional[str] = None
    ) -> Tuple[str, bool]:
        """
        确定 Go 类型和是否应该添加 omitempty 标签

        Args:
            field_type: 字段类型
            go_type: 初始 Go 类型
            required: 是否必需
            default_value: 默认值

        Returns:
            (最终Go类型, 是否需要omitempty)
        """
        # 如果有默认值或不是必需，添加 omitempty
        use_omitempty = not required or (default_value is not None and default_value.strip() != "")

        # 对于指针类型（非基础类型），使用指针以支持 nil
        if field_type in [FieldType.OBJECT, FieldType.ARRAY]:
            # Object 和 Array 使用指针以支持 nil 表示未设置
            return f"*{go_type}", use_omitempty

        return go_type, use_omitempty
