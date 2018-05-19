package com.ds.canopus.service;

import com.ds.canopus.domain.Token;
import com.ds.canopus.repository.TokenRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.time.LocalDate;

@Service
public class TokenServiceImpl implements TokenService {

    @Autowired
    private TokenRepository tokenRepository;

    @Override
    public Token findByToken(String token) {
        Token foundToken = tokenRepository.findTokenByToken(token);
        if (foundToken != null) {
            long delta = Duration.between(foundToken.getCreatedDt().toLocalDate().atStartOfDay(),
                    LocalDate.now().atStartOfDay()).toDays();
            if (delta <= 1 && delta >= 0) {
                return foundToken;
            }
        }
        return null;
    }

}
