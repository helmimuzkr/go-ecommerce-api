package service

import (
	"e-commerce-api/feature/product"
	"e-commerce-api/helper"
	"errors"
	"log"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
)

type productService struct {
	qry product.ProductData
	vld *validator.Validate
	cld *cloudinary.Cloudinary
}

func New(d product.ProductData, v *validator.Validate, cld *cloudinary.Cloudinary) product.ProductService {
	return &productService{
		qry: d,
		vld: v,
		cld: cld,
	}
}

func (ps *productService) Add(token interface{}, newProduct product.Core, fileHeader *multipart.FileHeader) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("token tidak valid")
	}

	if err := ps.vld.Struct(&newProduct); err != nil {
		msg := helper.ValidationErrorHandle(err)
		return errors.New(msg)
	}

	secureURL, err := helper.UploadFile(fileHeader, ps.cld)
	if err != nil {
		log.Println(err)
		var msg string
		if strings.Contains(err.Error(), "kesalahan input") {
			msg = err.Error()
		} else {
			msg = "gagal upload gambar karena kesalahan pada sistem server"
		}
		return errors.New(msg)
	}
	newProduct.Image = secureURL

	newProduct.Stock = 1

	if err := ps.qry.Add(uint(userID), newProduct); err != nil {
		return errors.New("terjadi kesalahan pada sistem server")
	}

	return nil
}

func (ps *productService) GetAll(page int) (map[string]interface{}, []product.Core, error) {
	// Total record
	totalRecord, err := ps.qry.CountProduct()
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("terjadi kesalahan pada sistem server")
	}
	if totalRecord < 1 {
		log.Println("total record kurang dari 1")
		return nil, nil, errors.New("data tidak ditemukan")
	}
	// Limit
	limit := 10
	// Total pages
	totalPage := totalRecord / limit
	if page >= totalPage {
		page = totalPage
	}
	if page < 2 {
		page = 1
	}
	// Calculate offset
	offset := (page - 1) * limit

	pagination := make(map[string]interface{})
	pagination["page"] = page
	pagination["limit"] = limit
	pagination["offset"] = offset
	pagination["totalRecord"] = totalRecord
	pagination["totalPage"] = totalPage

	res, err := ps.qry.GetAll(limit, offset)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("terjadi kesalahan pada sistem server")
	}

	return pagination, res, nil
}

func (ps *productService) GetByID(productID uint) (product.Core, error) {
	res, err := ps.qry.GetByID(productID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return product.Core{}, errors.New(msg)
	}

	return res, nil
}

func (ps *productService) Update(token interface{}, productID uint, updateProduct product.Core, fileHeader *multipart.FileHeader) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("token tidak valid")
	}

	if err := ps.vld.Struct(&updateProduct); err != nil {
		log.Println(err)
		msg := helper.ValidationErrorHandle(err)
		return errors.New(msg)
	}

	if fileHeader == nil {
		if err := ps.qry.Update(uint(userID), productID, updateProduct); err != nil {
			log.Println(err)
			msg := ""
			if strings.Contains(err.Error(), "tidak ditemukan") {
				msg = err.Error()
			} else {
				msg = "terjadi kesalahan pada sistem server"
			}
			return errors.New(msg)
		}

		return nil
	}

	// Proses update dan delete file
	// Ambil url avatar sebelumnya untuk dilakukan penghapusan file
	res, err := ps.qry.GetByID(productID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}
	secureURL, err := helper.UploadFile(fileHeader, ps.cld)
	if err != nil {
		log.Println(err)
		var msg string
		if strings.Contains(err.Error(), "kesalahan input") {
			msg = err.Error()
		} else {
			msg = "gagal upload gambar karena kesalahan pada sistem server"
		}
		return errors.New(msg)
	}
	updateProduct.Image = secureURL

	if err := ps.qry.Update(uint(userID), productID, updateProduct); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "tidak ditemukan") {
			msg = err.Error()
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}

	if res.Image != "" {
		publicID := helper.GetPublicID(res.Avatar)
		if err := helper.DestroyFile(publicID, ps.cld); err != nil {
			log.Println("destroy file", err)
			return errors.New("failed to destroy image")
		}
	}

	return nil
}

func (ps *productService) Delete(token interface{}, productID uint) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("token tidak valid")
	}

	// Proses update dan delete file
	// Ambil url avatar sebelumnya untuk dilakukan penghapusan file
	res, err := ps.qry.GetByID(productID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}

	// Delete data pada database
	if err := ps.qry.Delete(uint(userID), productID); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "tidak ditemukan") {
			msg = err.Error()
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}

	if res.Avatar != "" {
		// Delete file image
		publicID := helper.GetPublicID(res.Avatar)
		if err := helper.DestroyFile(publicID, ps.cld); err != nil {
			log.Println("destroy file", err)
			return errors.New("failed to destroy image")
		}
	}

	return nil
}
