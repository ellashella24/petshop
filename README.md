<div id="top"></div>

# Petshop

<!-- PROJECT LOGO -->
<div align="center">
  <a href="https://github.com/ellashella24/petshop">
    <img src="images/logo.png" alt="Logo" width="180" height="180">
  </a>

  <h3 align="center">Petshop</h3>


  <p align="center">
    A RESTful API for Petshop App
    <br />
    <div id = "other-software-design"></div>
    <a href="https://whimsical.com/petshop-RtwdxfQTB8e72AY681qRBj">Wireflow</a>
    ·
    <a href="https://lucid.app/lucidchart/7a103c4c-9aac-44e1-896f-b9bca4a8dc34/edit?invitationId=inv_c1d4dfe4-840d-4ce2-80fb-18845710d049">Use Case</a>
    ·
    <a href="https://lucid.app/lucidchart/8877930f-c2e0-4c85-b233-7c93dd82a306/edit?invitationId=inv_bccb8e74-b801-4419-ba8a-cc078d80f84a">Flowchart</a>
    ·
    <a href="https://lucid.app/lucidchart/1e35b3e5-1de0-40a7-ba26-db0e85763fda/edit?invitationId=inv_5eb6ec4e-9b6f-4db6-bb62-e18e2e922d95">ERD</a>
    ·
    <a href="https://app.swaggerhub.com/apis-docs/ellashella24/petshop/1.0.0">Open API</a>
  </p>
</div>
<br />

[![Contributors](https://img.shields.io/github/contributors/ellashella24/petshop.svg?style=for-the-badge)](https://github.com/ellashella24/petshop/graphs/contributors)

<!-- TABLE OF CONTENTS -->
# Table of Contents
1. [About the Project](#about-the-project)
2. [High Level Architecture](#high-level-architecture)
3. [Tech Stack](#tech-stack)
4. [Code Structure](#code-structure)
    - [Structuring](#structuring)
    - [Unit Test](#unit-test)
5. [How to Contrib](contribute.md)
6. [Contact](#contact)

<!-- ABOUT THE PROJECT -->
# About The Project
- An app that allow user to be a pet shop owner to sell their services and products or to be a customer to buy them. 
- Pet shop owners will be helped to market the products and services so that they can be easily reached by customers and they will be helped to get the products and services that are needed by their pets easily.
- Build with Golang, Echo Framework, MySQL adn GORM for manage repository, Xuri Excelize for Export List Product Selling to Excel, FTP to store Image Product to server, Xendit API for Payment Gateway, Deploy the project on [Okteto](https://ellashella24.cloud.okteto.net).

# High Level Architecture

HLA design for this project shown in the picture below

<img src="images/HLA-rev-3.jpeg" alt="hla" width="800" height="462" >

<br />

- Follow this link to see the other software design : <a href="#other-software-design">Other Software Design</a>

<p align="right">(<a href="#top">back to top</a>)</p>

# Tech Stack
## RESTful-API
- [Go](https://go.dev/)
- [Echo Framework](https://echo.labstack.com/)
- [MySQL](https://www.mysql.com/)
- [GORM](https://gorm.io/index.html)
- [FTP](https://github.com/jlaffaye/ftp)
- [SMTP](https://github.com/xhit/go-simple-mail)
- [Xuri Excelize](https://xuri.me/excelize/)
- [Xendit](https://www.xendit.co/id/?utm_source=google&utm_medium=cpc&utm_campaign=BKWS-Exact-ID-ID&utm_content=payment-gateway&utm_term=xendit)

## Deployment
- [Docker](https://www.docker.com/)
- [Okteto](https://www.okteto.com/)
- [Kubernetes](https://kubernetes.io/)

Follow the link to see deployment flow of this projetct : [Deployment Flow](https://lucid.app/lucidchart/cc522f9a-b238-4ad7-a546-84c70486c3fb/edit?invitationId=inv_3f370427-82c1-4217-a866-056b0021a77c)

## Collaboration 
- [Trello](https://trello.com/)
- [Github](https://github.com/)

<p align="right">(<a href="#top">back to top</a>)</p>

# Code Structure
This project use Layered Architure to organized each components into spesific function  

## Structuring
  ```sh
    petshop
    ├── config                        
    │     └──config.go                # Contains list of configuration of the project
    ├── constants                     
    │     └──constants.go             # Contains list constant variable
    ├── delivery                      # Contains list of component for handle request dan response
    │     └──common                   # Contains list of http request format based on the result from controller 
    │     │   ├── common.go           # Contains list of http request format
    │     │   └── http_responses.go   
    │     └──controller               # Contains list of component that receive the request and return a response
    │     │   ├── user
    │     │   ├── formatter_req.go    # Contains list of request format for each function on the controller
    │     │   ├── formatter_res.go    # Contains list of response format for each function on the controller
    │     │   ├── user_test.go        # Contains list of function for test each function on the controller
    │     │   └── users.go            # Contains list of controller for each entity    
    │     └──routes  
    │         └── routes.go           # Contains list of route to access each function on controller  
    ├── entity                        # Contains model all entity
    │     └── user.go                 # Contains model for spesific entity
    ├── repository                    # Contains list of functions that process the request and stores it in database
    │     ├── user_test.go            # Contains list of function for test each function on the repository
    │     └── users.go                # Contains list of repository for each entity
    ├── service                       # Contains list of function to access other components outside the project
    │     └── ftp.go                  # Containts list of function to manage files upload to the ftp
    │     └── export-excel.go         # Containts list of function to stores response data into the excel files
    ├── utils                         # Contains list of function to config each type of database
    │     └── mysqldriver.go          # Contains list of function to config MySQL type database
    ├── .env                          # Contains list of environment variable to run the project 
    ├── .gitignore                    # Contains list of directory/file name that will igonored when push project
    ├── go.mod                  
    ├── go.sum                  
    ├── main.go                       # Contains list of component that need to be executed first to run the app
    └── README.md    
  ```

## Unit Test
Coverage result on all functions is 99.2% which the most functions have reached 100% coverage. Coverage result for each function shown in the picture below

<img src="images/coverage-result-ver-2.jpg" alt="coverage-result">

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTACT -->
# Contact
* Naufal Aammar Hibatullah - [Github](https://github.com/nflhibatullah) · [LinkedIn](https://www.linkedin.com/in/naufal-hibatullah-441a58222/)
* Niendhitta Tamia Lassela - [Github](https://github.com/ellashella24) · [LinkedIn](https://www.linkedin.com/in/ntlassela/)

<p align="right">(<a href="#top">back to top</a>)</p>