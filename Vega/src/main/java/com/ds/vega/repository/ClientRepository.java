package com.ds.vega.repository;

import com.ds.vega.domain.Client;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface ClientRepository
        extends CrudRepository<Client, String> {
        //extends MongoRepository<Client, String> {
    List<Client> findAll();
    void deleteById(String id);
}
