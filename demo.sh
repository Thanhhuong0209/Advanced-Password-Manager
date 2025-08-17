#!/bin/bash

# Demo script for Advanced Password Manager
# This script demonstrates the main features of the password manager

echo " Advanced Password Manager Demo"
echo "=================================="
echo

# Check if the application is built
if [ ! -f "build/password-manager" ]; then
    echo "Building the application first..."
    make build
    echo
fi

# Set master password for demo (in real usage, this would be interactive)
export MASTER_PASSWORD="demo123"

echo "1.  Database Statistics"
echo "------------------------"
echo "master123" | ./build/password-manager stats
echo

echo "2. 🔑 Generate Strong Passwords"
echo "-------------------------------"
echo "master123" | ./build/password-manager generate --length 16 --uppercase --lowercase --numbers --symbols
echo
echo "master123" | ./build/password-manager generate --length 20 --uppercase --lowercase --numbers --symbols --no-repeating
echo
echo "master123" | ./build/password-manager generate --length 12 --uppercase --numbers
echo

echo "3. Save Passwords"
echo "-------------------"
echo "master123" | ./build/password-manager save gmail --username "user@gmail.com" --password "MySecurePass123!" --url "https://gmail.com" --notes "Personal email account" --tags "email,personal,google"
echo
echo "master123" | ./build/password-manager save github --username "developer" --password "DevPass456!" --url "https://github.com" --notes "GitHub account for development" --tags "development,git,code"
echo
echo "master123" | ./build/password-manager save bank --username "john.doe" --password "BankSecure789!" --url "https://mybank.com" --notes "Online banking account" --tags "banking,finance,personal"
echo

echo "4.  List All Passwords"
echo "------------------------"
echo "master123" | ./build/password-manager list
echo

echo "5. 🔍 Search Passwords"
echo "---------------------"
echo "master123" | ./build/password-manager search gmail
echo
echo "master123" | ./build/password-manager search personal
echo

echo "6. 📖 Retrieve Specific Password"
echo "-------------------------------"
echo "master123" | ./build/password-manager get gmail
echo

echo "7. 📊 Password Strength Analysis"
echo "-------------------------------"
echo "master123" | ./build/password-manager analyze "WeakPass123"
echo
echo "master123" | ./build/password-manager analyze "MyVerySecurePassword@2024!"
echo

echo "8. 📈 Final Database Statistics"
echo "------------------------------"
echo "master123" | ./build/password-manager stats
echo

echo "🎉 Demo Complete!"
echo "=================="
echo
echo "What you've seen:"
echo "✅ Password generation with various options"
echo "✅ Secure password storage with encryption"
echo "✅ Password retrieval and search"
echo "✅ Password strength analysis"
echo "✅ Database management"
echo
echo "Try these commands yourself:"
echo "  ./build/password-manager help"
echo "  ./build/password-manager generate --length 25"
echo "  ./build/password-manager save myaccount --username myuser --password mypass"
echo "  ./build/password-manager get myaccount"
echo
echo "For more information, see README.md"
