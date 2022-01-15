package conv

import "fmt"

type Inch float64
type Centimeter float64

func (i Inch) String() string {
	return fmt.Sprintf("%g in", i)
}

func (c Centimeter) String() string {
	return fmt.Sprintf("%g cm", c)
}

func CmToIn(cm Centimeter) Inch {
	return Inch(cm * 0.393701)
}

func InToCm(in Inch) Centimeter {
	return Centimeter(in / 0.393701)
}
