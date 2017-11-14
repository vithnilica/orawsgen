package orawsgen;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.Reader;
import java.io.StringReader;
import java.io.Writer;
import java.math.BigDecimal;
import java.sql.Blob;
import java.sql.CallableStatement;
import java.sql.Clob;
import java.sql.Connection;
import java.sql.SQLException;
import java.sql.SQLXML;
import java.util.Date;

import oracle.sql.OPAQUE;
import oracle.xdb.XMLType;


public class DefaultTypesMethods {

	//https://docs.oracle.com/cd/A97335_02/apps.102/a83724/basic3.htm

	//NUMBER
	public static void setNumberParam(Connection con, CallableStatement cs, String parameterName, BigDecimal data) throws SQLException {
		cs.setBigDecimal(parameterName, data);
	}
	public static void registerNumberOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		cs.registerOutParameter(parameterName, java.sql.Types.NUMERIC);
	}
	public static BigDecimal getNumberOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		return cs.getBigDecimal(parameterName);
	}
	public static Object getDbObjectNumber(java.sql.Connection con, BigDecimal o) throws java.sql.SQLException {
		if(o==null)return null;
		return o;
	}
	public static BigDecimal getWsObjectNumber(java.sql.Connection con, Object d) throws java.sql.SQLException {
		if(d==null)return null;
		return (BigDecimal) d;
	}

	//VARCHAR2
	public static void setVarchar2Param(Connection con, CallableStatement cs, String parameterName, String data) throws SQLException {
		cs.setString(parameterName, data);
	}
	public static void registerVarchar2OutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		cs.registerOutParameter(parameterName, java.sql.Types.VARCHAR);
	}
	public static String getVarchar2OutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		return cs.getString(parameterName);
	}
	public static Object getDbObjectVarchar2(java.sql.Connection con, String o) throws java.sql.SQLException {
		if(o==null)return null;
		return o;
	}
	public static String getWsObjectVarchar2(java.sql.Connection con, Object d) throws java.sql.SQLException {
		if(d==null)return null;
		return (String) d;
	}


	//DATE
	public static void setDateParam(Connection con, CallableStatement cs, String parameterName, Date data) throws SQLException {
		cs.setTimestamp(parameterName, DBUtils.convertDate2SqlTs(data));
	}
	public static void registerDateOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		cs.registerOutParameter(parameterName, java.sql.Types.TIMESTAMP);
	}
	public static Date getDateOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		return DBUtils.convertSqlTs2Date(cs.getTimestamp(parameterName));
	}
	public static Object getDbObjectDate(java.sql.Connection con, Date o) throws java.sql.SQLException {
		if(o==null)return null;
		return DBUtils.convertDate2SqlTs(o);
	}
	public static Date getWsObjectDate(java.sql.Connection con, Object d) throws java.sql.SQLException {
		if(d==null)return null;
		return DBUtils.convertSqlTs2Date((java.sql.Timestamp)d);
	}
	/*
	public static void setDateParam(Connection con, CallableStatement cs, String parameterName, XMLGregorianCalendar data) throws SQLException {
		cs.setTimestamp(parameterName, DBUtils.convertXmlCal2SqlTs(data));
	}
	public static void registerDateOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		cs.registerOutParameter(parameterName, java.sql.Types.TIMESTAMP);
	}
	public static XMLGregorianCalendar getDateOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		return DBUtils.convertSqlTs2XmlCal(cs.getTimestamp(parameterName));
	}
	public static Object getDbObjectDate(java.sql.Connection con, XMLGregorianCalendar o) throws java.sql.SQLException {
		if(o==null)return null;
		return DBUtils.convertXmlCal2SqlTs(o);
	}
	public static XMLGregorianCalendar getWsObjectDate(java.sql.Connection con, Object d) throws java.sql.SQLException {
		if(d==null)return null;
		return DBUtils.convertSqlTs2XmlCal((java.sql.Timestamp)d);
	}
	 */
	
	
	private static SQLXML castSafe2SQLXML(Object o)throws SQLException {
		if (o==null)return null;
		if (o instanceof SQLXML){
			return (SQLXML)o;
		}else if(o instanceof OPAQUE) {
			return XMLType.createXML((OPAQUE)o);
		} else {
			throw new SQLException("objekt nelze prevest na SQLXML "+o.getClass().getName());
		}
	}
	
	//XMLTYPE
	public static void setXmlTypeParam(Connection con, CallableStatement cs, String parameterName, XmlAny data) throws SQLException {
		if(data==null || data.any==null) {
			cs.setNull(parameterName, java.sql.Types.SQLXML,"XMLTYPE");
			//cs.setNull(parameterName, oracle.jdbc.OracleTypes.OPAQUE, "SYS.XMLTYPE");
		}else {
			cs.setSQLXML(parameterName, DBUtils.convertXmlAny2SqlXml(con, data));
		}
		//cs.setSQLXML(parameterName, DBUtils.convertXmlAny2SqlXml(con, data));
	}
	public static void registerXmlTypeOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		cs.registerOutParameter(parameterName, java.sql.Types.SQLXML);
	}
	public static XmlAny getXmlTypeOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		return DBUtils.convertSqlXml2XmlAny(cs.getSQLXML(parameterName));
	}
	public static Object getDbObjectXmlType(java.sql.Connection con, XmlAny o) throws java.sql.SQLException {
		if(o==null || o.any==null)return null;
		return DBUtils.convertXmlAny2SqlXml(con, o);
	}
	public static XmlAny getWsObjectXmlType(java.sql.Connection con, Object d) throws java.sql.SQLException {
		if(d==null)return null;
		return DBUtils.convertSqlXml2XmlAny(castSafe2SQLXML(d));
	}

	
	
	//CLOB
	public static void setClobParam(Connection con, CallableStatement cs, String parameterName, String data) throws SQLException {
		cs.setClob(parameterName, new StringReader(data));
	}
	public static void registerClobOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		cs.registerOutParameter(parameterName, java.sql.Types.CLOB);
	}
	public static String getClobOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		try {
			Clob c=cs.getClob(parameterName);
			if(c==null)return null;
			Reader r=c.getCharacterStream();
		    char[] arr = new char[8 * 1024];
		    StringBuilder buffer = new StringBuilder();
		    int numCharsRead;
		    while ((numCharsRead = r.read(arr, 0, arr.length)) != -1) {
		        buffer.append(arr, 0, numCharsRead);
		    }
		    r.close();
			return buffer.toString();	
		}catch (Exception e) {
			throw new SQLException(e);
		}

	}
	public static Object getDbObjectClob(java.sql.Connection con, String o) throws java.sql.SQLException {
		try {
			if(o==null)return null;
			Clob c=con.createClob();
			Writer w=c.setCharacterStream(1);
			w.write(o);
			w.close();
			return c;
		}catch (Exception e) {
			throw new SQLException(e);
		}
	}
	public static String getWsObjectClob(java.sql.Connection con, Object d) throws java.sql.SQLException {
		if(d==null)return null;
		try {
			Clob c=(Clob)d;
			Reader r=c.getCharacterStream();
		    char[] arr = new char[8 * 1024];
		    StringBuilder buffer = new StringBuilder();
		    int numCharsRead;
		    while ((numCharsRead = r.read(arr, 0, arr.length)) != -1) {
		        buffer.append(arr, 0, numCharsRead);
		    }
		    r.close();
			return buffer.toString();	
		}catch (Exception e) {
			throw new SQLException(e);
		}
	}
	
	
	//BLOB
	public static void setBlobParam(Connection con, CallableStatement cs, String parameterName, byte[] data) throws SQLException {
		cs.setBlob(parameterName, new ByteArrayInputStream(data));
	}
	public static void registerBlobOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		cs.registerOutParameter(parameterName, java.sql.Types.BLOB);
	}
	public static byte[] getBlobOutParam(Connection con, CallableStatement cs, String parameterName) throws SQLException {
		try {
			Blob b=cs.getBlob(parameterName);
			if(b==null)return null;
			InputStream is = b.getBinaryStream();
			ByteArrayOutputStream buffer = new ByteArrayOutputStream();
			int read;
			byte[] buf = new byte[8 * 1024];
			while ((read = is.read(buf, 0, buf.length)) != -1) {
				buffer.write(buf, 0, read);
			}
			return buffer.toByteArray();
		}catch (Exception e) {
			throw new SQLException(e);
		}
	}
	public static Object getDbObjectBlob(java.sql.Connection con, byte[] o) throws java.sql.SQLException {
		if(o==null)return null;
		try {
			Blob b=con.createBlob();
			OutputStream os=b.setBinaryStream(1);
			os.write(o);
			os.close();
			return b;
		}catch (Exception e) {
			throw new SQLException(e);
		}
	}
	public static byte[] getWsObjectBlob(java.sql.Connection con, Object d) throws java.sql.SQLException {
		if(d==null)return null;
		try {
			Blob b=(Blob)d;
			InputStream is = b.getBinaryStream();
			ByteArrayOutputStream buffer = new ByteArrayOutputStream();
			int read;
			byte[] buf = new byte[8 * 1024];
			while ((read = is.read(buf, 0, buf.length)) != -1) {
				buffer.write(buf, 0, read);
			}
			return buffer.toByteArray();
		}catch (Exception e) {
			throw new SQLException(e);
		}
	}

}
