package com.ds.vega.service;

import com.ds.vega.domain.Token;
import com.ds.vega.repository.TokenRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.time.LocalDate;
import java.time.ZoneId;
import java.util.Date;
import java.util.List;

@Service
public class TokenServiceImpl implements TokenService {

    @Autowired
    private TokenRepository tokenRepository;

    @Override
    public Token getTokenByTokenValue(String token) {
        List<Token> tokens = tokenRepository.findAll();
        Token foundToken = tokens.stream()
                .filter(token1 -> token.equals(token1.getToken()))
                .findAny()
                .orElse(null);

        if (foundToken != null) {
            Date createdDt = foundToken.getCreatedDt();
            LocalDate date = createdDt.toInstant().atZone(ZoneId.systemDefault()).toLocalDate();
            long delta = Duration.between(date.atStartOfDay(),
                    LocalDate.now().atStartOfDay()).toDays();
            if (delta <= 1 && delta >= 0) {
                return foundToken;
            }
        }
        return null;
    }

    @Override
    public List<Token> getAllTokens() {
        return tokenRepository.findAll();
    }

    @Override
    public Token insertToken(Token token) {
        return tokenRepository.save(token);
    }

    @Override
    public void deleteTokenById(String id) {
        tokenRepository.deleteById(id);
    }
}
