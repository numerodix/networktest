package main

//import "fmt"
import "testing"


func Test_Toolbox_haveTool(t *testing.T) {
    var tb = Toolbox{
        osName: "linux",
    }

    assertTrue(t, tb.haveTool("ls"), "ls not found")
    assertFalse(t, tb.haveTool("kjasdhfkjqwhjqa"), "kjasdhfkjqwhjqa found")
}
