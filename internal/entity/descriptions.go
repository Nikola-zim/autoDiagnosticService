package entity

type Classes struct {
	Classes map[int]string
}

const (
	numberOfClasses = 20
	accumulator     = "Accumulator: упало напряжение в бортовой сети, часто такая проблема связана с отсутствием заряда аккумуляторной батареи от генератора или проблемами с самим аккумулятором\n"
	abs             = "ABS: неполадки в антиблокировочной системе, в данный момент эта система не работает\n"
	oil             = "Oil: красная масленка свидетельствует о падении уровня масла в двигателе автомобиля.\n"
	engine          = "Engine: значок двигателя, он информирует о наличии ошибок двигателя и неисправности его электронных систем.\n"
	overheat        = "Overheat: значок охлаждающей жидкости, сообщает о повышенной температуре в системе охлаждения двигателя\n"
	brake           = "Brake: значок тормозной системы, говорит о неисправности тормозной системы\n"
	diselElectro    = "Disel/electro: сигнальная лампа предпускового подогрева дизельного двигателя или неисправности электронных систем.\n"
	airbag          = "Airbag: значок подушки безопасности, говорит о возникшей неисправности в системе пассивной безопасности, и в случае ДТП воздушные подушки не сработают\n"
	tyrePressure    = "Tyre Pressure – значок низкого давления в шинах. Предупреждает о падении давления воздуха ниже нормы в одном или нескольких шинах.\n"
	esp             = "ESP: значок может или периодически загораться или же гореть постоянно. Лампочка с такой надписью оповещает о проблемах системы стабилизации.\n"
	steering        = "PowerSteering: неисправность в электро- или гидроусилителе руля.\n"
)

func NewClasses() *Classes {
	classes := make(map[int]string, numberOfClasses)
	classes[0] = accumulator
	classes[1] = oil
	classes[2] = engine
	classes[3] = overheat
	classes[4] = brake
	classes[5] = diselElectro
	classes[6] = airbag
	classes[7] = tyrePressure
	classes[8] = esp
	classes[9] = abs
	classes[10] = ""
	classes[11] = steering
	classes[12] = ""
	classes[13] = ""
	classes[14] = ""
	classes[15] = ""
	classes[16] = ""
	classes[17] = ""
	classes[18] = ""
	classes[19] = ""
	classes[20] = ""
	classes[21] = ""
	classes[22] = ""
	return &Classes{
		Classes: classes,
	}
}
