package main

import "fmt"

type Animal struct{ Name string }

func (a Animal) Speak() string    { return a.Name + " makes a sound" }
func (a Animal) Describe() string { return "I am " + a.Name }

type Dog struct {
	Animal // embedded — promotes Speak() and Describe()
	Breed  string
}

// Override promoted method
func (d Dog) Speak() string { return d.Name + " says: Woof!" }

type ServiceDog struct {
	Dog  // double embedding
	Role string
}

func (s ServiceDog) Describe() string {
	return fmt.Sprintf("%s is a %s service dog", s.Name, s.Role)
}

// Interface satisfied by embedding
type Speaker interface{ Speak() string }

func makeSpeak(s Speaker) { fmt.Println(s.Speak()) }

func main() {
	a := Animal{Name: "Cat"}
	d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
	sd := ServiceDog{Dog: Dog{Animal: Animal{Name: "Buddy"}}, Role: "Guide"}

	fmt.Println(a.Speak())
	fmt.Println(d.Speak())        // overridden
	fmt.Println(d.Describe())     // promoted from Animal
	fmt.Println(sd.Speak())       // promoted from Dog
	fmt.Println(sd.Describe())    // overridden

	makeSpeak(d)  // Dog satisfies Speaker via Speak()
	makeSpeak(sd) // ServiceDog promotes Dog.Speak()
}
