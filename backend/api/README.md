### rustful api

#### /sign_in

**headers**

method: post

Content-Type: application/x-www-form-urlencoded 

**body**

```
name     string
password string
```

#### /sign_up

**headers**

```
method: post
Content-Type: application/x-www-form-urlencoded 
```

**body**

```
name                      string
password                  string
new_password_confirmation string
```

#### /renew_password

**headers**

```
method: post
Content-Type: application/x-www-form-urlencoded 
```

**body**

```
name                      string
password                  string
new_password              string
new_password_confirmation string
```

#### /sign_out

method: post

**limit**

> authorization required


#### /takeout

**headers**

```
method: post
Content-Type: application/json
```

**body**

```
{
	authuser: string
	address:  string
	tag:      string
	items: [{
		name: string
		price: uint
		num: int
	}]
}
```

**limit**

> authorization required

#### /order

**headers**

```
method: post
Content-Type: application/json
```

**body**

```
{
	order_id: string
}
```

#### /order_list

**headers**

```
method: post
Content-Type: application/json
```

**body**

```
{
	start: iso 8601 time format string
	end: iso 8601 time format string
}
```


