## This is the repo for the FSF mobile backend.

`go/` is for our Go files. It houses the code to interact with CiviCRM and TrustCommerce: accepting payments and handling CAS login. `See go/README.md` for more detailed documentation on the Go files.
`rails/` is for our Rails server. It houses endpoints and database schemas for the petition page and the newsfeed.  
`third_party/` contains third-party libraries that we integrate. Currently, it just contains the TCLink library in C.

To start, 

1. $ cd ../third_party/tclink/
2. $ ./configure

This script configures the software on your specific system. It makes sure all of the dependencies for the rest of the build and install process are available, and finds out whatever it needs to know to use those dependencies.

3. $ make
4. $ cd ../../go/
5. $ make
6. $ ./cas_server -tcpasswd <PASSWORD HERE> -tcuser <USER NAME HERE> -adminkey <ADMIN KEY HERE> -sitekey <SITE KEY HERE>