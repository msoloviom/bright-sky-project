package com.ds.vega;

import com.ds.vega.service.ClientService;
import com.ds.vega.service.TokenService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.data.redis.connection.jedis.JedisConnectionFactory;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.GenericToStringSerializer;


@SpringBootApplication
public class VegaApplication {

	@Autowired
	private static ClientService clientService;

	@Autowired
	private TokenService tokenService;

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
                .append("token", "lkdfndcvkashgc")
                .append("clientId", "1")
        		.append("createdDt", new Date());
		collection.insertOne(client);
		//collection.drop();
		mongoClient.close();*/

		SpringApplication.run(VegaApplication.class, args);
	}

	@Bean
	JedisConnectionFactory jedisConnectionFactory() {

		JedisConnectionFactory connectionFactory = new JedisConnectionFactory();
/*		connectionFactory.setHostName("192.168.0.82");
		connectionFactory.setPort(6380);*/
		return connectionFactory;
	}

	@Bean

	public RedisTemplate<String, Object> redisTemplate() {

		final RedisTemplate<String, Object> template = new RedisTemplate<String, Object>();
		template.setConnectionFactory(jedisConnectionFactory());
		template.setValueSerializer(new GenericToStringSerializer<Object>(Object.class));
		return template;

	}

}
