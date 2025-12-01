"""
数据类型和模型定义
用于表示解析后的 API 信息、字段信息等
"""

from dataclasses import dataclass, field
from typing import List, Optional, Dict, Any
from enum import Enum


class FieldType(Enum):
    """字段类型枚举"""
    # 基础类型
    INT64 = "int64"
    INT32 = "int32"
    INT = "int"
    UINT64 = "uint64"
    UINT32 = "uint32"
    UINT = "uint"
    STRING = "string"
    BOOL = "bool"
    FLOAT64 = "float64"
    FLOAT32 = "float32"

    # 复合类型
    OBJECT = "object"
    ARRAY = "array"

    # OneBot 特殊类型
    MESSAGE = "message"  # 可以是字符串或消息段数组

    # 未知类型
    UNKNOWN = "unknown"


class MessageTypeVariant(Enum):
    """Message 类型的变体"""
    STRING = "string"          # 纯文本字符串
    ARRAY = "array"            # 消息段数组


@dataclass
class Field:
    """字段定义"""
    name: str                          # 字段名（如：user_id）
    go_name: str                       # Go 字段名（如：UserID）
    data_type: str                     # 原始数据类型字符串（如：number (int64)）
    field_type: FieldType              # 解析后的字段类型
    go_type: str                       # Go 类型（如：int64, string）
    description: str                   # 字段描述
    required: bool = True              # 是否必需
    default_value: Optional[str] = None  # 默认值
    possible_values: List[str] = field(default_factory=list)  # 可能的值
    is_optional: bool = False          # 是否通过 omitempty 标记为可选

    # 仅用于 MESSAGE 类型
    message_variants: List[MessageTypeVariant] = field(default_factory=list)

    # 仅用于 ARRAY 类型
    element_type: Optional[str] = None
    element_field_type: Optional[FieldType] = None

    # 仅用于 OBJECT 类型
    nested_fields: List['Field'] = field(default_factory=list)


@dataclass
class APIModel:
    """API 请求/响应模型"""
    name: str                          # 模型名称（如：SendPrivateMsgReq）
    fields: List[Field]                # 字段列表
    description: str = ""              # 模型描述
    doc_link: str = ""                 # 文档链接


@dataclass
class APIDefinition:
    """API 定义"""
    api_name: str                      # API 名称（如：send_private_msg）
    description: str                   # API 描述
    request_model: APIModel            # 请求模型
    response_model: APIModel           # 响应模型
    doc_link: str = ""                 # 文档链接


@dataclass
class EventModel:
    """事件模型"""
    name: str                          # 事件名称（如：PrivateMessageEvent）
    fields: List[Field]                # 字段列表
    event_type: str = ""               # 事件类型（如：message, notice）
    sub_type: str = ""                 # 事件子类型（如：friend, group）
    description: str = ""              # 事件描述
    doc_link: str = ""                 # 文档链接


@dataclass
class MessageSegment:
    """消息段定义"""
    segment_type: str                  # 消息段类型（如：text, image）
    fields: List[Field]                # 字段列表
    description: str = ""              # 消息段描述
    can_send: bool = True              # 是否可发送
    can_receive: bool = True           # 是否可接收
    doc_link: str = ""                 # 文档链接
