#!/usr/bin/env python3
"""
OneBot 11 Go SDK 代码生成器主程序

用法:
    python main.py [选项]

选项:
    --input-dir      输入 Markdown 文档目录 (默认: ../api)
    --output-dir     输出 Go 代码目录 (默认: ../output)
    --package        Go 包名 (默认: onebot)
"""

import argparse
import sys
from pathlib import Path

from go_generator import GoCodeGenerator
from markdown_parser import MarkdownParser
from schema import *


def main():
    parser = argparse.ArgumentParser(description="OneBot 11 Go SDK 代码生成器")

    parser.add_argument(
        "--input-dir",
        type=str,
        default="../../",
        help="输入 Markdown 文档根目录（包含 api 和 event 子目录）",
    )
    parser.add_argument(
        "--output-dir",
        type=str,
        default="../models",
        help="输出 Go 代码目录",
    )
    parser.add_argument(
        "--package",
        type=str,
        default="models",
        help="Go 包名",
    )

    args = parser.parse_args()

    # 转换为绝对路径
    input_dir = Path(args.input_dir).resolve()
    output_dir = Path(args.output_dir).resolve()

    print(f"📖 输入目录: {input_dir}")
    print(f"📝 输出目录: {output_dir}")
    print(f"📦 Go 包名: {args.package}")

    # 检查输入目录
    if not input_dir.exists():
        print(f"❌ 错误: 输入目录不存在: {input_dir}")
        return 1

    # 创建输出目录
    output_dir.mkdir(parents=True, exist_ok=True)

    # 初始化解析器和生成器
    markdown_parser = MarkdownParser()
    go_generator = GoCodeGenerator(package_name=args.package)

    try:
        # ========== 解析 API 文档 ==========
        api_dir = input_dir / "api"
        api_file = api_dir / "public.md"

        if api_file.exists():
            print(f"\n🔍 解析 API 文档: {api_file}")
            apis = markdown_parser.parse_api_file(str(api_file))
            print(f"✅ 成功解析 {len(apis)} 个 API")

            # 生成 API Go 代码
            print(f"\n⚙️  生成 API Go 代码...")
            api_code = go_generator.generate_all_apis(apis)

            # 写入输出文件
            api_output_file = output_dir / "api.go"
            with open(api_output_file, "w", encoding="utf-8") as f:
                f.write(api_code)

            print(f"✅ 成功生成: {api_output_file}")
        else:
            print(f"⚠️  警告: 未找到 API 文档: {api_file}")
            apis = []

        # ========== 解析事件文档 ==========
        event_dir = input_dir / "event"
        event_files = [
            event_dir / "message.md",
            event_dir / "notice.md",
            event_dir / "request.md",
            event_dir / "meta.md",
        ]

        all_events = []
        for event_file in event_files:
            if event_file.exists():
                print(f"\n🔍 解析事件文档: {event_file}")
                events = markdown_parser.parse_event_file(str(event_file))
                print(f"✅ 成功解析 {len(events)} 个事件")
                all_events.extend(events)
            else:
                print(f"⚠️  警告: 未找到事件文档: {event_file}")

        if all_events:
            # 生成事件 Go 代码
            print(f"\n⚙️  生成事件 Go 代码...")
            event_code = go_generator.generate_all_events(all_events)

            # 写入输出文件
            event_output_file = output_dir / "event.go"
            with open(event_output_file, "w", encoding="utf-8") as f:
                f.write(event_code)

            print(f"✅ 成功生成: {event_output_file}")

        # ========== 解析消息段文档 ==========
        message_dir = input_dir / "message"
        message_file = message_dir / "segment.md"

        all_segments = [] # type: List[MessageSegment]
        if message_file.exists():
            print(f"\n🔍 解析消息段文档: {message_file}")
            segments = markdown_parser.parse_message_segment_file(str(message_file))
            print(f"✅ 成功解析 {len(segments)} 个消息段")
            all_segments.extend(segments)
        else:
            print(f"⚠️  警告: 未找到消息段文档: {message_file}")

        if all_segments:
            # 生成消息段 Go 代码
            print(f"\n⚙️  生成消息段 Go 代码...")
            message_code = go_generator.generate_all_message_segments(all_segments)

            # 写入输出文件
            message_output_file = output_dir / "message.go"
            with open(message_output_file, "w", encoding="utf-8") as f:
                f.write(message_code)

            print(f"✅ 成功生成: {message_output_file}")

        # 统计信息
        print(f"\n📊 生成统计:")
        print(f"  - API 数量: {len(apis)}")
        print(f"  - 事件数量: {len(all_events)}")
        print(f"  - 消息段数量: {len(all_segments)}")
        if apis:
            print(f"  - API 输出文件: {output_dir / 'api.go'}")
        if all_events:
            print(f"  - 事件输出文件: {output_dir / 'event.go'}")
        if all_segments:
            print(f"  - 消息段输出文件: {output_dir / 'message.go'}")

        return 0

    except Exception as e:
        print(f"❌ 错误: {e}")
        import traceback
        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
