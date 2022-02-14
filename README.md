<div id="top"></div>
<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/ellashella24/petshop">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Petshop</h3>

  <p align="center">
    An RESTful API Application for Petshop
    <br />
    <a href="https://github.com/ellashella24/petshop"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/ellashella24/petshop">View Demo</a>
    ·
    <a href="https://github.com/ellashella24/petshop/issues">Report Bug</a>
    ·
    <a href="https://github.com/ellashella24/petshop/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#restful-api">RESTful-API</a></li>
      <ul>
        <li><a href="#user-related">User Related</a></li>
        <li><a href="#pet-related">Pet Related</a></li>
        <li><a href="#store-related">Store Related</a></li>
        <li><a href="#product-related">Product Related</a></li>
        <li><a href="#cart-related">Cart Related</a></li>
        <li><a href="#transaction-related">Transaction Related</a></li>
        <li><a href="#category-related">Category Related</a></li>
        <li><a href="#city-related">City Related</a></li>
      </ul>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>
<br/>

<!-- ABOUT THE PROJECT -->
## About The Project

A RESTful API App for Petshop. Build with Golang, Echo Framework, MySQL adn GORM for manage repository, Deploy on Okteto

<p align="right">(<a href="#top">back to top</a>)</p>

### Built With

* [Go](https://go.dev/)
* [Echo Framework](https://echo.labstack.com/)
* [MySQL](https://www.mysql.com/)
* [GORM](https://gorm.io/index.html)


<p align="right">(<a href="#top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

This is how to use the project

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/ellashella24/petshop.git
   ```
2. Run the program
    ```
    go run main.go
    ```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- RESTFUL API -->

## RESTful-API

### Open-endpoints
Open endpoints require no Authentication. 

- Register : `POST /user/register`
- Login : `POST /user/login`

### Endpoints that require Authentication
Closed endpoints require a valid Token to be included in the header of the request. A Token can be acquired from the Login view above.

### User Related
Each endpoint manipulates or displays information related to the User whose Token is provided with the request:

- Get user own profile : `GET /user/profile`
- Update user own profile data : `PUT /user/profile`
- Delete user own account : `DELETE /user`

### Pet Related
Each endpoint manipulates or displays information related to the Pet whose Token is provided with the request:

- Create pet by user: `POST /user/pet`
- Get all pet by user : `GET /user/pets`
- Get pet profile : `GET /user/pet/profile/:id`
- Update pet profile : `PUT /user/pet/profile/:id`
- Delete pet by user : `DELETE /user/pet/:id`
- Get grooming status : `GET /user/grooming_status/pet/:id`
- Update 'SELESAI' status on grooming pet status : `PUT /user/grooming_status/pet/:id`

### Store Related
Each endpoint manipulates or displays information related to the Store whose Token is provided with the request:

- Create store by user: `POST /user/store`
- Get all store by user : `GET /user/stores`
- Get store profile : `GET /user/store/profile/:id`
- Update store profile : `PUT /user/store/profile/:id`
- Delete store by user : `DELETE /user/store/:id`
- Get grooming status of pet : `GET /store/grooming_status/pet/:id`
- Update grooming status of pet : `PUT /store/grooming_status/pet/:id` 
- Export product sales list to Excel : `GET /export/transactions/store/:id`

### Product Related
Each endpoint manipulates or displays information related to the Product whose Token is provided with the request:

- Create product by store: `POST /product`
- Get all product : `GET /product`
- Get product by id : `GET /product/:id`
- Get all product by store: `GET /product?store=`
- Update product by id : `PUT /product/:id`
- Delete product by id : `DELETE /product/:id`
- Get update stock history of product : `GET /stock/product/:id`

### Cart Related
Each endpoint manipulates or displays information related to the Cart whose Token is provided with the request:

- Create cart by user: `POST /cart`
- Checkout product from cart : `POST /cart/checkout`
- Get all cart : `GET /cart`
- Delete cart by id : `DELETE /cart/:id`

### Transaction Related
Each endpoint manipulates or displays information related to the Transaction whose Token is provided with the request:

- Create new transaction: `POST /transaction`
- Get all transaction by user : `GET /transaction/user/:id`
- Get all transaction by store : `GET /transaction/store/:id`
- Get callback : `POST /callback`

## Endpoints that require Admin Role
The endpoint below requires checking that the currently logged in user role is admin

### Category Related
Each endpoint manipulates or displays information related to the Category whose Token is provided with the request:

- Create category : `POST /category`
- Get all category : `GET /categories`
- Get category by id : `GET /category/profile/:id`
- Update category : `PUT /category/profile/:id`
- Delete category : `DELETE /category/:id`

### City Related
Each endpoint manipulates or displays information related to the City whose Token and the role is admin that provided with the request:

- Create city : `POST /city`
- Get all city : `GET /cities`
- Get city by id : `GET /city/profile/:id`
- Update city : `PUT /city/profile/:id`
- Delete city : `DELETE /city/:id`

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTACT -->
## Contact
* [Naufal Aamar Hibatullah](https://github.com/nflhibatullah)
* [Niendhitta Tamia Lassela](https://github.com/ellashella24)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Layered Architecture](https://github.com/jackthepanda96/docker-be5)
* [Readme Template](https://github.com/othneildrew/Best-README-Template)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/othneildrew/Best-README-Template.svg?style=for-the-badge
[contributors-url]: https://github.com/ellashella24/petshop/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/othneildrew/Best-README-Template.svg?style=for-the-badge
[forks-url]:  https://github.com/ellashella24/petshop/network/members
[stars-shield]: https://img.shields.io/github/stars/othneildrew/Best-README-Template.svg?style=for-the-badge
[stars-url]:  https://github.com/ellashella24/petshop/stargazers
[issues-shield]: https://img.shields.io/github/issues/othneildrew/Best-README-Template.svg?style=for-the-badge
[issues-url]:  https://github.com/ellashella24/petshop/issues
[linkedin-shield-ela]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url-ela]: https://www.linkedin.com/in/ntlassela
[linkedin-shield-naufal]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url-naufal]: https://www.linkedin.com/in/naufal-hibatullah-441a58222
