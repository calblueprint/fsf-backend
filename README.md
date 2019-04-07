# Backend for FSF mobile app.

### Folder structure

`go/` is for our Go files. It houses the code to interact with CiviCRM and TrustCommerce: accepting payments and handling CAS login. `See go/README.md` for more detailed documentation on the Go files.

`rails/` is for our Rails server. It houses endpoints and database schemas for the petition page and the newsfeed.  

`third_party/` contains third-party libraries that we integrate. Currently, it just contains the TCLink library in C.

### How to run the go backend (`go/`)

The go backend depends on tclink in the `third_party` folder. You need to first build tclink before you build the go folder.

To build tclink:

```
cd third_party/tclink/
./configure
make
```

Once tclink is built, you can build the go backend using make. If you have not, you need to install go.

```
cd go/
make
```

Run the go backend binary:
```
./cas_server -tcpasswd <PASSWORD HERE> -tcuser <USER NAME HERE> -adminkey <ADMIN KEY HERE> -sitekey <SITE KEY HERE>
```


