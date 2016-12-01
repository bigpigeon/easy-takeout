### rustful api

#### /sign_in

method: post


Content-Type: application/x-www-form-urlencoded
Content-Type: json

**data**

> name     string

> password string


#### /sign_up

method: post


Content-Type: application/x-www-form-urlencoded 
Content-Type: json

**data**

> name                      string

> password                  string

> new_password_confirmation string

#### /renew_password

method: post


Content-Type: application/x-www-form-urlencoded
Content-Type: json

**data**

> name                      string

> password                  string

> new_password              string

> new_password_confirmation string

#### /sign_out

method: post

