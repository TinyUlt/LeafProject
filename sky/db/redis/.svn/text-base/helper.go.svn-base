
package redis

import (
    "net"
    "time"
    "bufio"
    "io/ioutil"
    "sky/jsonhelper"
)

// Dial connects to the given Redis server.
func Dial(node jsonhelper.JSONObject) *Client {
    var (
        bufLen  int
        nt      string
    )
    if bufLen = node.GetAsInt("buff_size"); bufLen == 0 {
        bufLen = 4096
    }
    if nt = node.GetAsString("net"); nt == "" {
        nt = "tcp"
    }
    c, err := dial(nt, node.GetAsString("addr"), 
        bufLen, time.Duration(node.GetAsInt("timeout")))
    if err != nil {
        panic(err)
    }
    if pwd := node.GetAsString("password"); len(pwd) > 0 {
        if err = c.Cmd("auth", pwd).Err; err != nil {
            panic(err)
        }
    }
    if which := node.GetAsInt("db"); which > 0 {
        if err = c.Cmd("select", which).Err; err != nil {
            panic(err)
        }
    }
    return c
}

func Close(c *Client) {
    if c != nil {
        c.Close()
    }
}


// Dial connects to the given Redis server.
func LoadScript(c *Client, path string) (sha string) {
    var b   []byte
    var err error
    if b, err = ioutil.ReadFile(path); err != nil {
        panic(err)
    }
    if sha, err = c.Cmd("SCRIPT", "LOAD", string(b)).String(); err != nil {
        panic(err)
    }
    return
}

// Dial connects to the given Redis server.
func FastDial(network, addr string) (*Client, error) {
    return dial(network, addr, 4096, time.Duration(0))
}

// Dial connects to the given Redis server with the given timeout.
func dial(network, addr string, bufLen int, timeout time.Duration) (*Client, error) {
    // establish a connection
    //log.Trace("Dial redis %s", addr)
    conn, err := net.Dial(network, addr)
    if err != nil {
        return nil, err
    }
    //log.Trace("Dial redis %s succeed", addr)
    
    c := new(Client)
    c.conn = conn
    c.timeout = timeout
    c.reader = bufio.NewReaderSize(conn, bufLen)

    return c, nil
}