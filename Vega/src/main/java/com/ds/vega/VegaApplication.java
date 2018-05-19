package com.ds.vega;

import com.ds.vega.repository.ClientRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;


@SpringBootApplication
public class VegaApplication {

	@Autowired
	private ClientRepository clientRepository;

	public static void main(String[] args) {
		/*MongoClient mongoClient = new MongoClient();
		MongoDatabase database = mongoClient.getDatabase("vega");
        Document client = new Document("name", "John Doe")
                //.append("name", "John Doe")
        		.append("cert", "iyfliyf8r86ru6r8udp8o6uod68");
		MongoCollection<Document> collection = database.getCollection("client");
		collection.insertOne(client);
		//collection.drop();
		mongoClient.close();*/

		/*MongoClient mongoClient = new MongoClient();
		MongoDatabase database = mongoClient.getDatabase("vega");
		MongoCollection<Document> collection = database.getCollection("token");
        Document client = new Document("_id", "1")
                .append("token", "iyfgraigaerggu")
                .append("clientId", "1")
        		.append("createdDt", new Date());
		collection.insertOne(client);
		//collection.drop();
		mongoClient.close();*/


		SpringApplication.run(VegaApplication.class, args);
	}
}
