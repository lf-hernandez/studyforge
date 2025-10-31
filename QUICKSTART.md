# StudyForge Quick Start Guide

Get StudyForge running in 5 minutes!

## Step 1: Install Go

If you don't have Go installed:

### macOS
```bash
brew install go
```

### Linux
```bash
# Download and install
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Windows
Download installer from: https://go.dev/dl/

## Step 2: Get a Hugging Face API Key

1. Go to https://huggingface.co/
2. Sign up for a free account
3. Go to Settings → Access Tokens
4. Create a new token (read access is enough)
5. Copy your token

## Step 3: Configure Environment

```bash
cd studyforge

# Copy environment file
cp .env.example .env

# Edit .env and add your API key
nano .env  # or use any text editor
```

Replace `your_api_key_here` with your actual Hugging Face API key.

## Step 4: Install Dependencies

```bash
go mod download
```

## Step 5: Run the Application

```bash
# Make the run script executable
chmod +x scripts/run.sh

# Run the application
./scripts/run.sh
```

Or run directly:
```bash
go run cmd/server/main.go
```

## Step 6: Open in Browser

Navigate to: **http://localhost:8080**

## What's Next?

1. Upload a PDF textbook (max 50MB)
2. Select a page range (e.g., 1-10)
3. Choose your academic level
4. Click "Generate Summary"
5. View your AI-generated study notes!

## Troubleshooting

**Problem**: "Go command not found"
- **Solution**: Install Go (see Step 1) and restart your terminal

**Problem**: "API key error"
- **Solution**: Verify your Hugging Face API key is correct in `.env`

**Problem**: "Port 8080 already in use"
- **Solution**: Change `SERVER_PORT` in `.env` to a different port (e.g., 8081)

**Problem**: "Failed to process PDF"
- **Solution**:
  - Ensure the PDF is not corrupted
  - Try a different PDF file
  - Check file size is under 50MB
  - Some PDFs are image-based and need OCR (coming in Milestone 3)

## Need Help?

- Check the full README.md for detailed documentation
- Review markdown/PROJECT_OUTLINE.md for architecture details
- See markdown/PROJECT_GUIDELINES.md for development standards

---

Happy studying! 📚✨
