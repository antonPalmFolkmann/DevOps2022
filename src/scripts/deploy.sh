docker kill minitwit
docker rm minitwit
docker pull antonfolkmann/minitwit
docker run -d -p 8081:8081 --name minitwit antonfolkmann/minitwit