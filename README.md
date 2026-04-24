# SSO Access Project

## Setup Documentation

Follow these steps to get the project running locally.

### Backend (Golang) Setup

1. **Clone the repository**
2. **Install dependencies**:
   ```
   cd golang
   go mod tidy
   ```
3. **Ensure Services are Running**: Make sure MySQL and Redis are already running on your machine.
4. **Environment Variables**: Copy `.env.example` to `.env` and fill in the actual values.
   ```bash
   cp .env.example .env
   ```
5. **Run Migrations**:
   ```bash
   make migrate
   ```
6. **Run Database Seeding**:
   ```bash
   make seed
   ```
7. **Start the Server** (using Air for hot reloading):
   ```bash
   air
   ```
   
`default username: super-admin`
`default password: very-secret`

### Frontend (Next.js) Setup

1. **Install dependencies**:
   ```bash
   cd next
   npm install
   ```
2. **Configure environment variables**:
   Copy the example environment file and configure your local variables:
   ```bash
   cp .env.example .env
   ```
3. **Run the development server**:
   Start the Next.js development server:
   ```bash
   npm run dev
   ```
   
The application will be available at [http://localhost:3000](http://localhost:3000).