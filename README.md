# Pancake autoloader for Docter and Commander

An Arknights and Girls' Frontline news watcher for docter and commander

## How to use

Before using autoloader, you need to build a QQ message sender, [mirai-cpp-messagesender](
https://github.com/D-0000000000/mirai-cpp-messagesender)

### Build autoloader

First get source code
```bash
$ git clone https://github.com/D-0000000000/autoloader.git
```

Specific QQ message sender for function `consume` in `autoloader.go`.
If you want to use a new QQ message sender, you can rewrite function `consume`. 

Then build autoloader.

```bash
$ cd autoloader
$ go build
```

## Usage

```bash
$ autoloader
```

Once autoloader get new message, message will send to QQ by message sender

## Extra Infomation

Actually I don't know anything about golang. I just modified a project from [hguandl](https://github.com/hguandl/) I don't know if there is any problem. But it did work. If there is any problems on code I'm out.

Watch list refer to `config.yaml`.

## Credit

This project is based on [hguandl](https://github.com/hguandl/)'s project [dr-feeder](https://github.com/hguandl/dr-feeder).  
I just modified some configuration and call a externel program to send message.
