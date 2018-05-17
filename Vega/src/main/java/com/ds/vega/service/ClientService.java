package com.ds.vega.service;

import com.ds.vega.domain.Client;

import java.util.List;

public interface ClientService {

    Client getClientById(String id);
    List<Client> getAllClients();
    Client insertClient(Client client);

}
