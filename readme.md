# DevPulse 🚀

![CodeRabbit Reviews](https://img.shields.io/coderabbit/prs/github/sharmatanikesh/devplus?utm_source=oss&utm_medium=github&utm_campaign=sharmatanikesh%2Fdevplus&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)
![Vercel](https://img.shields.io/badge/Deployed%20on-Vercel-000000?style=flat&logo=vercel&logoColor=white)
![Kestra](https://img.shields.io/badge/Powered%20by-Kestra-4A90E2?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEyIDJMMiAxMkwxMiAyMkwyMiAxMkwxMiAyWiIgZmlsbD0id2hpdGUiLz4KPC9zdmc+)

**DevPulse** is an AI-powered engineering intelligence platform that provides automated PR reviews, intelligent release notes generation, and comprehensive engineering metrics. Gain actionable insights into your development workflow and improve team productivity with data-driven decisions.

## ✨ Key Features

- 🤖 **AI-Powered PR Analysis** - Automatically review pull requests with intelligent feedback and recommendations
- 📦 **Smart Release Notes** - Generate comprehensive, well-formatted release notes and changelogs
- 📊 **Engineering Metrics Dashboard** - Track PR lead time, review latency, merge frequency, and team performance
- 🔍 **Impact Analysis** - Identify architectural hotspots and code coupling issues
- 🎯 **Release Risk Assessment** - Evaluate deployment risks before releases
- 🔗 **Seamless GitHub Integration** - OAuth authentication and real-time webhook support
- 📈 **Personal Developer Metrics** - Individual contribution tracking and insights

## 🏗️ Architecture

DevPulse is built with a modern, scalable architecture:

### Frontend

- **Framework**: Next.js 16 with App Router
- **UI**: ShadCN UI + Tailwind CSS
- **Language**: TypeScript
- **Deployment**: Vercel
- **State Management**: Redux Toolkit

### Backend

- **Language**: Go (Golang)
- **Framework**: Gorilla Mux
- **Database**: PostgreSQL (GORM)
- **Deployment**: Google Cloud Platform (GCP)
- **Authentication**: GitHub OAuth 2.0
- **AI Integration**: Kestra workflows for intelligent analysis

## 🚀 Quick Start

### Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js** 20+ and npm/yarn
- **Go** 1.20+
- **PostgreSQL** 13+
- **Docker & Docker Compose** (optional, for containerized setup)

### GitHub OAuth Setup

1. Create a GitHub OAuth App:
   - Go to **Settings** > **Developer settings** > **OAuth Apps** > **New OAuth App**
   - Set **Authorization callback URL** to `http://localhost:8080/api/v1/auth/github/callback`
   - Note your **Client ID** and **Client Secret**

### Local Development Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/sharmatanikesh/devplus.git
   cd devplus
   ```

2. **Start Kestra with Docker Compose**

   Kestra is the AI workflow engine that powers intelligent analysis:

   ```bash
   cd kestra
   docker run -d \
     --name kestra \
     -p 8080:8080 \
     -v $(pwd)/storage:/app/storage \
     kestra/kestra:latest
   ```

   Or build from the Dockerfile:

   ```bash
   cd kestra
   docker build -t kestra-local .
   docker run -d --name kestra -p 8080:8080 kestra-local
   ```

   This will start Kestra at http://localhost:8080

3. **Backend Setup**

   ```bash
   cd backend
   cp .env.example .env
   # Edit .env with your configuration
   make run
   ```

   The backend will start at http://localhost:8081

   See [Backend README](./backend/README.md) for detailed instructions.

4. **Frontend Setup**

   ```bash
   cd frontend
   npm install
   cp .env.example .env.local
   # Edit .env.local with your configuration
   npm run dev
   ```

   See [Frontend README](./frontend/README.md) for detailed instructions.

5. **Access the Application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8081
   - Kestra UI: http://localhost:8080

## 🧪 Code Quality

This project uses **[CodeRabbit AI](https://coderabbit.ai)** for automated code reviews on every pull request. CodeRabbit provides:

- Intelligent PR reviews with actionable feedback
- Security vulnerability detection
- Code quality suggestions
- Best practice recommendations

All pull requests are automatically reviewed by CodeRabbit to maintain high code quality standards.

## 📊 Project Structure

```
devplus/
├── frontend/              # Next.js frontend application
│   ├── app/              # Next.js App Router pages
│   │   ├── (dashboard)/  # Protected dashboard routes
│   │   ├── login/        # Authentication page
│   │   ├── privacy/      # Privacy policy
│   │   └── terms/        # Terms of service
│   ├── components/       # React components
│   │   ├── dashboard/    # Dashboard components
│   │   ├── repositories/ # Repository components
│   │   └── ui/           # ShadCN UI components
│   ├── hooks/            # React hooks
│   ├── lib/              # Utilities and API client
│   │   ├── store/        # Redux store
│   │   ├── types/        # TypeScript types
│   │   └── utils/        # Utility functions
│   └── public/           # Static assets
├── backend/              # Go backend service
│   ├── cmd/              # Application entry point
│   │   └── server/       # Main server
│   ├── internal/         # Internal packages
│   │   ├── config/       # Configuration
│   │   ├── controllers/  # HTTP handlers
│   │   │   └── rest/     # REST controllers
│   │   ├── services/     # Business logic
│   │   │   ├── ai/       # AI service factory
│   │   │   ├── auth_service/
│   │   │   └── github_service/
│   │   ├── models/       # Data models
│   │   ├── repositories/ # Data access layer
│   │   ├── middleware/   # HTTP middleware
│   │   ├── router/       # Route definitions
│   │   ├── db/           # Database connection
│   │   └── migrations/   # SQL migrations
│   ├── pkg/              # Public packages
│   │   ├── logger/       # Logging utilities
│   │   └── utils/        # Utility functions
│   ├── workflows/        # Kestra workflow YAML files
│   ├── Dockerfile        # Backend container
│   ├── docker-compose.yaml
│   └── Makefile          # Build commands
├── kestra/               # Kestra workflow engine
│   ├── Dockerfile        # Kestra container
│   └── README.md         # Kestra setup guide
└── readme.md             # This file
```

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

All PRs will be automatically reviewed by CodeRabbit AI before human review.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Tanikesh Sharma** - [@sharmatanikesh](https://github.com/sharmatanikesh)
- **Shobhit Tomer** - [@ShobhitTomer](https://github.com/ShobhitTomer)

## 🙏 Acknowledgments

- [CodeRabbit](https://coderabbit.ai) for AI-powered code reviews
- [Vercel](https://vercel.com) for seamless deployment
- [Kestra](https://kestra.io) for workflow automation

## 📞 Support

For issues, questions, or suggestions:

- 🐛 [Open an issue](https://github.com/sharmatanikesh/devplus/issues)
- 💬 [Start a discussion](https://github.com/sharmatanikesh/devplus/discussions)

---

**Made with ❤️ from ShaiTan, for developers**
