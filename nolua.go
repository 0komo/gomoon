//go:build !(lua54 || lua53 || lua52 || lua51)

package gomoon

// i wish Go has compile-time error 😔
const errorMessage = `
gomoon requires one of these build tags to be defined:
  - lua54
  - lua53
  - lua52
  - lua51

make sure to include that when building
`

func init() {
	panic(errorMessage)
}
