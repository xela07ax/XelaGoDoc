package encodingStdout

import (
	"fmt"
	"github.com/xela07ax/toolsXela/tp"
	"testing"
)

func TestConvert(t *testing.T) {
	f866, _ := tp.OpenReadFile("oute866.txt")
	fUTF, _ := tp.OpenReadFile("outeUTF.txt")
	fmt.Printf("--->cp866:%s\n--->utf:%s\n",f866,fUTF)
	fNewUtf := Convert(41,f866)
	fmt.Printf("--->newUtf:%s\n",fNewUtf)
    if len(fUTF) != len(fNewUtf) {
        t.Error("Кодирование не удалось", "")
    }
}