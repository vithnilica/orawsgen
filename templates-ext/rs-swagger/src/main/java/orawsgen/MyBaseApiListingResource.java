package orawsgen;

import javax.servlet.ServletConfig;
import javax.servlet.ServletContext;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.HttpHeaders;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.UriInfo;

import com.fasterxml.jackson.core.JsonProcessingException;

import io.swagger.jaxrs.listing.BaseApiListingResource;
import io.swagger.models.Swagger;
import io.swagger.util.Json;
import io.swagger.util.Yaml;

//HOTFIX nahrada za io.swagger.jaxrs.listing.BaseApiListingResource
public class MyBaseApiListingResource extends BaseApiListingResource{

	@Override
    protected Response getListingJsonResponse(
            Application app,
            ServletContext servletContext,
            ServletConfig servletConfig,
            HttpHeaders headers,
            UriInfo uriInfo) throws JsonProcessingException {
        Swagger swagger = process(app, servletContext, servletConfig, headers, uriInfo);
        
        if (swagger != null) {
            //HOTFIX z nakeho duvodu original celou cestu do konfiguraku nedava, ale aby ui fungovalo, tak je to potreba
            swagger.setBasePath(uriInfo.getBaseUri().getPath());
            
            return Response.ok().entity(Json.mapper().writeValueAsString(swagger)).type(MediaType.APPLICATION_JSON_TYPE).build();
        } else {
            return Response.status(404).build();
        }
    }
	
	@Override
    protected Response getListingYamlResponse(
            Application app,
            ServletContext servletContext,
            ServletConfig servletConfig,
            HttpHeaders headers,
            UriInfo uriInfo) {
        Swagger swagger = process(app, servletContext, servletConfig, headers, uriInfo);
                
        try {
            if (swagger != null) {
                //HOTFIX z nakeho duvodu original celou cestu do konfiguraku nedava, ale aby ui fungovalo, tak je to potreba
                swagger.setBasePath(uriInfo.getBaseUri().getPath());
                
                String yaml = Yaml.mapper().writeValueAsString(swagger);
                StringBuilder b = new StringBuilder();
                String[] parts = yaml.split("\n");
                for (String part : parts) {
                    b.append(part);
                    b.append("\n");
                }
                return Response.ok().entity(b.toString()).type("application/yaml").build();
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        return Response.status(404).build();
    }
}
