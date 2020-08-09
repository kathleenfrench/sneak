package helpers

import "fmt"

// Spacer is a simple util for adding space between terminal messages
func Spacer() {
	fmt.Println("")
}

// func Sudo(cmd string) (string, error) {
// 			// sudo check
// 			a := exec.Command("sudo", "ls")
// 			a.Stderr = os.Stderr
// 			a.Stdin = os.Stdin
// 			out, err := a.Output()
// 			if err != nil {
// 				gui.ExitWithError(err)
// 			}

// 			fmt.Println(string(out))
// }
