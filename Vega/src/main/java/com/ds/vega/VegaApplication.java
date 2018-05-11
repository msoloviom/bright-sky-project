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
        Document client = new Document("_id", "2")
                .append("name", "Jessica Jones");
		MongoCollection<Document> collection = database.getCollection("client");
		collection.insertOne(client);
		//ArrayList<Client> results = collection.find().into(new ArrayList<>());
		mongoClient.close();*/


		SpringApplication.run(VegaApplication.class, args);
	}
}
