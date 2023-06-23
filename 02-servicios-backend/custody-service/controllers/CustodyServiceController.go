package controllers

import (
	"context"
	"errors"
	"regexp"

	pb "github.com/malarcon-79/microservices-lab/grpc-protos-go/system/custody"
	"github.com/malarcon-79/microservices-lab/orm-go/dao"
	"github.com/malarcon-79/microservices-lab/orm-go/model"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Controlador de servicio gRPC
type CustodyServiceController struct {
	logger *zap.SugaredLogger // Logger
	re     *regexp.Regexp     // Expresión regular para validar formato de períodos YYYY-MM
}

// Método a nivel de package, permite construir una instancia correcta del controlador de servicio gRPC
func NewCustodyServiceController() (CustodyServiceController, error) {
	_logger, _ := zap.NewProduction() // Generamos instancia de logger
	logger := _logger.Sugar()

	re, err := regexp.Compile(`^\d{4}\-(0?[1-9]|1[012])$`) // Expresión regular para validar períodos YYYY-MM
	if err != nil {
		return CustodyServiceController{}, err
	}

	instance := CustodyServiceController{
		logger: logger, // Asignamos el logger
		re:     re,     // Asignamos el RegExp precompilado
	}
	return instance, nil // Devolvemos la nueva instancia de este Struct y un puntero nulo para el error
}

func (c *CustodyServiceController) AddCustodyStock(ctx context.Context, msg *pb.CustodyAdd) (*pb.Empty, error) {
	// instancio para trabajar con la tabla Custody
	orm := dao.DB.Model(&model.Custody{})

	// Validaciones

	// valido que el periodo tenga info
	if len(msg.Period) == 0 {
		return nil, errors.New("Período nulo")
	}

	// valido el periodo con la expresión regular
	if err := c.re.MatchString(msg.Period); !err {
		c.logger.Error("formato de Período inválido", err)
		return nil, errors.New("formato de Período inválido")
	}

	// valido que el stock venga con info
	if len(msg.Stock) == 0 {
		return nil, errors.New("Stock inválido")
	}

	// valido que el client id venga con info
	if len(msg.ClientId) == 0 {
		return nil, errors.New("ID de cliente inválido")
	}

	// valido que el Quantity no sea menor a 0
	if msg.Quantity < 0 {
		return nil, errors.New("Cantidad debe ser mayor a cero")
	}

	// valido que el quantity sea in int32 por definición de BD
	if msg.Quantity != float64(int32(msg.Quantity)) {
		return nil, errors.New("Cantidad debe ser un int")
	}

	// creo el modelo de datos para almacenamiento
	custody := &model.Custody{
		Period:   msg.Period,
		Stock:    msg.Stock,
		ClientId: msg.ClientId,
		Quantity: int32(msg.Quantity),
		Market:   "",
		Price:    decimal.NewFromInt(0),
	}

	// Ejecutamos el INSERT y evaluamos si hubo errores
	if err := orm.Create(custody).Error; err != nil {
		c.logger.Error("no se pudo guardar la custodia nueva", err)
		return nil, errors.New("error al guardar la custodia nueva")
	}

	// Retornamos la respuesta Empty
	return &pb.Empty{}, nil
}

func (c *CustodyServiceController) ClosePeriod(ctx context.Context, msg *pb.CloseFilters) (*pb.Empty, error) {
	return nil, errors.New("no implementado")
}

func (c *CustodyServiceController) GetCustody(ctx context.Context, msg *pb.CustodyFilter) (*pb.Custodies, error) {
	// instancio ORM para trabajar con Custody
	orm := dao.DB.Model(&model.Custody{})

	// Arreglo de punteros a registros de tabla "Invoice"
	custodies := []*model.Custody{}

	// Creamos el filtro de búsqueda usando los campos del mismo modelo
	filter := &model.Custody{
		Period:   msg.Period,
		Stock:    msg.Stock,
		ClientId: msg.ClientId,
	}

	// Ejecutamos el SELECT con el filtro
	if err := orm.Find(&custodies, filter).Error; err != nil {
		c.logger.Errorf("no se pudo buscar custodias con filtros %v", filter, err)
		return nil, status.Errorf(codes.Internal, "no se pudo realizar query")
	}

	// Este será el mensaje de salida
	result := &pb.Custodies{}

	// iteramos el resultado de la query
	for _, item := range custodies {

		result.Items = append(result.Items, &pb.Custodies_Custody{
			Period:   item.Period,
			Stock:    item.Stock,
			ClientId: item.ClientId,
			Quantity: item.Quantity,
		})
	}

	// Retornamos la respuesta
	return result, nil
}
