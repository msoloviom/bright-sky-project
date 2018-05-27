package com.ds.vega.service;

import com.ds.vega.domain.Client;

import java.util.List;
import java.util.Optional;

public interface ClientService {

    Optional<Client> getClientById(String id);
    List<Client> getAllClients();
    Client insertClient(Client client);
    void deleteClientById(String id);
}
