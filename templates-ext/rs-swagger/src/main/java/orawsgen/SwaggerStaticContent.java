package orawsgen;
import java.io.InputStream;
import java.net.URI;

import javax.servlet.ServletContext;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.UriInfo;

@Path("")
public class SwaggerStaticContent {
    @Context
    ServletContext servletContext;
	@Context
    private UriInfo uriInfo;
    
    @Path("") 
    @GET
    public Response index() {
    	URI uri = uriInfo.getBaseUriBuilder().path("/index.html").build();
        return Response.temporaryRedirect(uri).build();
    }
    
    @Path("/{path: .+}")
    @GET
    public InputStream getFile(@PathParam("path") String path) {
    	System.out.println("getFile "+path);
    	if("".equals(path)||"/".equals(path))path="index.html";
    	InputStream is = servletContext.getResourceAsStream("/WEB-INF/classes/swagger/"+path);
    	return is;
    }
}

