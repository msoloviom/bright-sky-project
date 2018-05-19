package com.ds.vega.repository;

import com.ds.vega.domain.Token;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface TokenRepository extends MongoRepository<Token, String> {

    Token findTokenByToken(String token);
    List<Token> findAll();
    Token insert(Token token);
}
