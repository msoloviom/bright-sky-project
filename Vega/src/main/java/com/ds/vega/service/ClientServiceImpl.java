package com.ds.vega.service;

import com.ds.vega.domain.Client;
import com.ds.vega.repository.ClientRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ClientServiceImpl implements ClientService {

    @Autowired
    private ClientRepository clientRepository;

    @Override
    public Client getClientById(String id) {
        return clientRepository.findClientById(id);
    }

    @Override
    public List<Client> getAllClients() {
        return clientRepository.findAll();
    }

    @Override
    public Client insertClient(Client client) {
        return clientRepository.insert(client);
    }
}
