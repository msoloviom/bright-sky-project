package com.ds.vega.domain;

import org.springframework.data.mongodb.core.mapping.Document;

@Document(collection = "client")
public class Client {

    private String id;

    private String name;

    private String cert;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getCert() {
        return cert;
    }

    public void setCert(String cert) {
        this.cert = cert;
    }
}
