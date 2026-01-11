#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  OTP Email Configuration Setup${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ .env file created${NC}"
    echo ""
fi

# Check if SMTP credentials already exist
if grep -q "SMTP_EMAIL=" .env && grep -q "SMTP_PASSWORD=" .env; then
    echo -e "${YELLOW}SMTP credentials already exist in .env file.${NC}"
    read -p "Do you want to update them? (y/n): " update_choice
    if [ "$update_choice" != "y" ]; then
        echo -e "${GREEN}Keeping existing SMTP credentials.${NC}"
        exit 0
    fi
fi

echo -e "${YELLOW}Please provide your Gmail SMTP credentials:${NC}"
echo ""

# Get SMTP email
read -p "Gmail address: " smtp_email
if [ -z "$smtp_email" ]; then
    echo -e "${RED}Error: Email cannot be empty${NC}"
    exit 1
fi

# Get SMTP password
echo ""
echo -e "${YELLOW}Google App Password (the one provided: yefb iude pmjo askn):${NC}"
read -p "App Password: " smtp_password
if [ -z "$smtp_password" ]; then
    echo -e "${RED}Error: Password cannot be empty${NC}"
    exit 1
fi

# Add or update SMTP credentials in .env
if grep -q "SMTP_EMAIL=" .env; then
    # Update existing
    sed -i.bak "s|SMTP_EMAIL=.*|SMTP_EMAIL=$smtp_email|" .env
    sed -i.bak "s|SMTP_PASSWORD=.*|SMTP_PASSWORD=$smtp_password|" .env
    rm .env.bak 2>/dev/null
else
    # Add new
    echo "" >> .env
    echo "# Email Configuration (Gmail SMTP)" >> .env
    echo "SMTP_EMAIL=$smtp_email" >> .env
    echo "SMTP_PASSWORD=$smtp_password" >> .env
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ SMTP credentials configured successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Restart your backend server (go run main.go)"
echo "2. Test the signup flow at http://localhost:3000/signup"
echo "3. Check OTP_SETUP.md for detailed documentation"
echo ""
