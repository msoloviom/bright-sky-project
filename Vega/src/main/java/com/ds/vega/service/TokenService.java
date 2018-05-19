package com.ds.vega.service;

import com.ds.vega.domain.Token;

import java.util.List;

public interface TokenService {

    Token getTokenByTokenValue(String token);

    List<Token> getAllTokens();

    Token insertToken(Token token);
}
