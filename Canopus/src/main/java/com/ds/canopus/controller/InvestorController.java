package com.ds.canopus.controller;

import com.ds.canopus.domain.Investor;
import com.ds.canopus.domain.Token;
import com.ds.canopus.repository.TokenRepository;
import com.ds.canopus.resource.InvestorResource;
import com.ds.canopus.resource.PostResource;
import com.ds.canopus.service.InvestorService;
import com.ds.canopus.service.TokenService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

import static org.springframework.util.CollectionUtils.isEmpty;

@RestController
@RequestMapping("/api")
public class InvestorController {

    private static final String SUCCESSFULLY_CREATED = "Investor was successfully created !";

    @Autowired
    private InvestorService investorService;

    @Autowired
    private TokenService tokenService;

    @Autowired
    private TokenRepository tokenRepository;

    @GetMapping("/investors")
    public HttpEntity<List<Investor>> getAllInvestors() {
        List<Investor> investors = investorService.getAllInvestors();
        if(isEmpty(investors)) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(investors, HttpStatus.OK);
        }
    }

    @GetMapping("/investors/{id}")
    public HttpEntity<InvestorResource> getInvestorById(@PathVariable Long id) {
        Investor investor = investorService.getInvestorById(id);
        System.out.println("\n*** This instance is picked for a call! " +
                "******************************************\n");

        if (investor == null) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            InvestorResource investorResource = new InvestorResource();
            investorResource.setId(investor.getId().toString());
            investorResource.setName(investor.getName());
            investorResource.setCert(investor.getCert());
            return new ResponseEntity<>(investorResource, HttpStatus.OK);
        }
    }

    @PostMapping("/investors")
    public HttpEntity<PostResource> saveInvestor(@RequestBody Investor investor) {
        Long investorId = investorService.saveInvestor(investor);
        return new ResponseEntity<>(new PostResource(investorId, SUCCESSFULLY_CREATED), HttpStatus.CREATED);
    }

    @GetMapping("/investors/auth/{token}")
    public HttpEntity<Long> getInvestorIdByToken(@PathVariable String token) {
        Token foundToken = tokenService.findByToken(token);
        if (foundToken == null) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(foundToken.getInvestorId(), HttpStatus.OK);
        }
    }

    @GetMapping("/investors/auth/tokens")
    public HttpEntity<List<Token>> getInvestorTokens() {

        List<Token> tokens = tokenRepository.findAll();

        if (isEmpty(tokens)) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(tokens, HttpStatus.OK);
        }
    }
}
