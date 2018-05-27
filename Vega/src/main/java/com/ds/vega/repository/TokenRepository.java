package com.ds.vega.repository;

import com.ds.vega.domain.Token;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface TokenRepository
        extends CrudRepository<Token, String> {
        //extends MongoRepository<Token, String> {

    Token findTokenByToken(String token);
    List<Token> findAll();
}
