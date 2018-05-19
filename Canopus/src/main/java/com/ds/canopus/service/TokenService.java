package com.ds.canopus.service;

import com.ds.canopus.domain.Token;

public interface TokenService {

    Token findByToken(String token);
}
