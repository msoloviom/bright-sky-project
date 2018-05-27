package com.ds.vega.service;

import com.ds.vega.domain.Client;
import com.ds.vega.repository.ClientRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class ClientServiceImpl implements ClientService {

    @Autowired
    private ClientRepository clientRepository;

    @Override
    public Optional<Client> getClientById(String id) {
        return clientRepository.findById(id);
    }

    @Override
    public List<Client> getAllClients() {
        return clientRepository.findAll();
    }

    @Override
    public Client insertClient(Client client) {
        return clientRepository.save(client);
    }

    @Override
    public void deleteClientById(String id) {
        clientRepository.deleteById(id);
    }

}
