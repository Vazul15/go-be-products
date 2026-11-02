package router

import (
	"products/internal/connector/restserver/endpoint"

	"github.com/gin-gonic/gin"
)

const (
	basePath        = "/api/products"
	webhook         = "/webhook"
	checkoutSession = "/checkout-session"
	getProduct      = "/:public_id"
	getVariant      = "/variants/:public_id"
)

type Router struct {
	getProduct           *endpoint.GetProduct
	listProducts         *endpoint.ListProducts
	createProduct        *endpoint.CreateProduct
	updateProduct        *endpoint.UpdateProduct
	deleteProduct        *endpoint.DeleteProduct
	createProductVariant *endpoint.CreateProductVariant
	listProductVariants  *endpoint.ListProductVariants
	getProductVariant    *endpoint.GetProductVariant
	updateProductVariant *endpoint.UpdateProductVariant
	deleteProductVariant *endpoint.DeleteProductVariant
}

// Register registers the routes and their authorization middleware.
func (r *Router) Register(rg gin.IRouter) {
	rg.GET(basePath, r.listProducts.Endpoint)
	rg.POST(basePath, r.createProduct.Endpoint)
	rg.POST(basePath+getProduct+"/variants", r.createProductVariant.Endpoint)
	rg.GET(basePath+getProduct+"/variants", r.listProductVariants.Endpoint)
	rg.GET(basePath+getVariant, r.getProductVariant.Endpoint)
	rg.PATCH(basePath+getVariant, r.updateProductVariant.Endpoint)
	rg.DELETE(basePath+getVariant, r.deleteProductVariant.Endpoint)
	rg.PATCH(basePath+getProduct, r.updateProduct.Endpoint)
	rg.DELETE(basePath+getProduct, r.deleteProduct.Endpoint)
	rg.GET(basePath+getProduct, r.getProduct.Endpoint)
}

// From older project as tempalte
// rg.POST(exportPath, r.exportAuth.Middleware, r.exportData.Endpoint)
