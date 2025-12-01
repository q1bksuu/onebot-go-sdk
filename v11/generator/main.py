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

import sys
import argparse
from pathlib import Path
from markdown_parser import MarkdownParser
from go_generator import GoCodeGenerator


def main():
    parser = argparse.ArgumentParser(description="OneBot 11 Go SDK 代码生成器")

    parser.add_argument(
        "--input-dir",
        type=str,
        default="../api",
        help="输入 Markdown 文档目录",
    )
    parser.add_argument(
        "--output-dir",
        type=str,
        default="../output",
        help="输出 Go 代码目录",
    )
    parser.add_argument(
        "--package",
        type=str,
        default="onebot",
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

    # 查找并解析 API 文档
    api_file = input_dir / "public.md"
    if not api_file.exists():
        print(f"❌ 错误: 找不到 API 文档: {api_file}")
        return 1

    print(f"\n🔍 解析 API 文档: {api_file}")

    try:
        apis = markdown_parser.parse_api_file(str(api_file))
        print(f"✅ 成功解析 {len(apis)} 个 API")

        # 生成 Go 代码
        print(f"\n⚙️  生成 Go 代码...")
        go_code = go_generator.generate_all_apis(apis)

        # 写入输出文件
        models_file = output_dir / "models.go"
        with open(models_file, "w", encoding="utf-8") as f:
            f.write(go_code)

        print(f"✅ 成功生成: {models_file}")

        # 统计信息
        print(f"\n📊 生成统计:")
        print(f"  - API 数量: {len(apis)}")
        print(f"  - 输出文件: {models_file}")

        return 0

    except Exception as e:
        print(f"❌ 错误: {e}")
        import traceback
        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
