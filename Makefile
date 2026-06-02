.PHONY: up down clean install

up: install
	@echo "🚀 Starting Agni services in LOCAL mode (SQLite + Local Redis)..."
	@if [ ! -f .env ]; then \
		echo "⚠️  .env file missing. Copying from .env.sample..."; \
		cp .env.sample .env; \
	fi
	@# Start Backend
	@ENV_MODE=local go run cmd/server/main.go > backend.log 2>&1 & echo $$! > .backend.pid
	@# Start In-App Server
	@ENV_MODE=local go run cmd/inapp/main.go > inapp.log 2>&1 & echo $$! > .inapp.pid
	@sleep 2
	@echo "🖥️  Starting Web Frontend..."
	@ENV_MODE=local cd web && npm run dev

install:
	@echo "📦 Checking and installing frontend dependencies..."
	@if [ ! -d web/node_modules ]; then \
		cd web && npm install; \
	fi

down:
	@echo "🛑 Stopping local Agni processes..."
	@if [ -f .backend.pid ]; then \
		kill -9 $$(cat .backend.pid) 2>/dev/null || true; \
		rm .backend.pid; \
	fi
	@if [ -f .inapp.pid ]; then \
		kill -9 $$(cat .inapp.pid) 2>/dev/null || true; \
		rm .inapp.pid; \
	fi
	@echo "✅ All local background processes stopped."

clean:
	@echo "🧹 Cleaning local state..."
	@rm -f agni.db backend.log inapp.log .backend.pid .inapp.pid
