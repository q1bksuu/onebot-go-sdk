.PHONY: help generate clean install test

help:
	@echo "OneBot 11 Go SDK Generator"
	@echo ""
	@echo "Available targets:"
	@echo "  generate    - Generate Go code from Markdown documentation"
	@echo "  clean       - Remove generated files"
	@echo "  install     - Install Python dependencies (if needed)"
	@echo "  test        - Run tests (if available)"
	@echo ""
	@echo "Examples:"
	@echo "  make generate"
	@echo "  make generate input-dir=../api output-dir=./output"

generate:
	@echo "🔍 Generating Go code from Markdown..."
	python3 generator/main.py $(if $(input-dir),--input-dir $(input-dir)) $(if $(output-dir),--output-dir $(output-dir))
	@echo "✅ Done!"

clean:
	@echo "🧹 Cleaning up..."
	rm -f output/models.go
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete
	@echo "✅ Cleaned!"

install:
	@echo "📦 Installing dependencies..."
	uv pip install -e .
	@echo "✅ Dependencies installed!"

test:
	@echo "🧪 Running tests..."
	@echo "⚠️  Tests not yet implemented"

.PHONY: docs

docs:
	@echo "📖 Opening README..."
	@open README.md || less README.md
