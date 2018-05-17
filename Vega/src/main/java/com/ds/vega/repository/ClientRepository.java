package com.ds.vega.repository;

import com.ds.vega.domain.Client;
import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.List;

public interface ClientRepository extends MongoRepository<Client, String> {
    Client findClientById(String id);
    List<Client> findAll();
    Client insert(Client client);
}
