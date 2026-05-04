//go:build ignore
// Remove the above line when implementing

package main

import "fmt"

type Animal struct{ Name string }

// TODO: Animal.Speak returns "<Name> makes a sound"
func (a Animal) Speak() string {
	// TODO: implement
	return ""
}

type Dog struct {
	Animal
	Breed string
}

// TODO: Dog.Speak overrides Animal.Speak, returns "<Name> says: Woof!"
func (d Dog) Speak() string {
	// TODO: implement
	return ""
}

type ServiceDog struct {
	Dog
	Role string
}

// TODO: ServiceDog.Describe returns "<Name> is a <Role> service dog"
func (s ServiceDog) Describe() string {
	// TODO: implement
	return ""
}

type Speaker interface{ Speak() string }

func makeSpeak(s Speaker) { fmt.Println(s.Speak()) }

func main() {
	a := Animal{Name: "Cat"}
	d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
	sd := ServiceDog{Dog: Dog{Animal: Animal{Name: "Buddy"}}, Role: "Guide"}

	fmt.Println(a.Speak())      // Expected: Cat makes a sound
	fmt.Println(d.Speak())      // Expected: Rex says: Woof!
	fmt.Println(d.Describe())   // Expected: I am Rex  (promoted from Animal)
	fmt.Println(sd.Describe())  // Expected: Buddy is a Guide service dog
	makeSpeak(d)                // Expected: Rex says: Woof!
	makeSpeak(sd)               // Expected: Buddy says: Woof! (promoted)
}
