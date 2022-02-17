# How to contribute

- Fork this repository

    ```sh
    $ git clone https://github.com/YOUR_USERNAME/petshop.git
    > Cloning into `petshop`...
    > remote: Counting objects: 10, done.
    > remote: Compressing objects: 100% (8/8), done.
    > remove: Total 10 (delta 1), reused 10 (delta 1)
    > Unpacking objects: 100% (10/10), done.
    ```

    ```sh
    cd petshop
    ```

- Important

    Always create new branch when develop something

    ```sh
    git checkout -b feature-name 
    ```

    ```sh
    git add .    
    ```

    ```sh
    git commit -m "feature description"
    ```

    ```sh
    $ git remote -v
    > origin  https://github.com/YOUR_USERNAME/petshop.git (fetch)
    > origin  https://github.com/YOUR_USERNAME/petshop.git (push)
    ```

    ```sh
    git remote add upstream https://github.com/ellashella24/petshop.git
    ```

    ```sh
    $ git remote -v
    > origin    https://github.com/YOUR_USERNAME/petshop.git (fetch)
    > origin    https://github.com/YOUR_USERNAME/petshop.git (push)
    > upstream  https://github.com/ellashella24/petshop.git (fetch)
    > upstream  https://github.com/ellashella24/petshop.git (push)
    ```

    ```sh
    git push -u origin feature-name    
    ```