package com.ds.canopus.repository;

import com.ds.canopus.domain.Token;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface TokenRepository extends CrudRepository<Token, Long> {

    Token findTokenByToken(String token);

    List<Token> findAll();

}
