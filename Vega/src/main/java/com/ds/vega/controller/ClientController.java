package com.ds.vega.controller;

import com.ds.vega.domain.Client;
import com.ds.vega.domain.Token;
import com.ds.vega.resource.PostResource;
import com.ds.vega.service.ClientService;
import com.ds.vega.service.TokenService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Date;
import java.util.List;

import static org.springframework.util.CollectionUtils.isEmpty;

@RestController
@RequestMapping("/api")
public class ClientController {

    private static final String SUCCESSFULLY_CREATED = "Client was successfully created";

    @Autowired
    private ClientService clientService;

    @Autowired
    private TokenService tokenService;

    @GetMapping(value = "/clients")
    public ResponseEntity<List<Client>> getAllClients() {
        List<Client> clients = clientService.getAllClients();
        if (isEmpty(clients)) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(clients, HttpStatus.OK);
        }
    }

    @GetMapping(value = "/clients/{id}")
    public ResponseEntity<Client> getClientById(@PathVariable("id") String id) {
        Client client = clientService.getClientById(id).orElse(null);
        if (client == null) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(client, HttpStatus.OK);
        }
    }

    @PostMapping(value = "/clients")
    public HttpEntity<PostResource> saveClient(@RequestBody Client clientRequest) {
        Client client = clientService.insertClient(clientRequest);
        return new ResponseEntity<>(new PostResource(client.getId(), SUCCESSFULLY_CREATED),
                HttpStatus.CREATED);

    }

    @DeleteMapping(value = "/clients/{id}")
    public HttpEntity<PostResource> deleteClient(@PathVariable String id) {
        clientService.deleteClientById(id);
        return new ResponseEntity<>(HttpStatus.OK);
    }


    @GetMapping(value = "/clients/auth/{token}")
    public HttpEntity<Long> getClientIdByToken(@PathVariable String token) {
        Token foundToken = tokenService.getTokenByTokenValue(token);
        if (foundToken == null) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(foundToken.getClientId(), HttpStatus.OK);
        }
    }

    @GetMapping(value = "/clients/auth/tokens")
    public HttpEntity<List<Token>> getClientTokens() {
        List<Token> tokens = tokenService.getAllTokens();
        if (isEmpty(tokens)) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(tokens, HttpStatus.OK);
        }
    }

    @PostMapping(value = "/clients/auth/tokens")
    public HttpEntity<PostResource> saveToken(@RequestBody Token tokenReq) {
        tokenReq.setCreatedDt(new Date());
        tokenService.insertToken(tokenReq);
        return new ResponseEntity<>(HttpStatus.CREATED);
    }

    @DeleteMapping(value = "/clients/auth/tokens/{id}")
    public HttpEntity<PostResource> saveToken(@PathVariable String id) {
        tokenService.deleteTokenById(id);
        return new ResponseEntity<>(HttpStatus.OK);
    }


}
