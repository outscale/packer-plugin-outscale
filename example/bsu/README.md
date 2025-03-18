# **Outscale Packer Examples**

This folder contains multiple examples demonstrating how to use **Packer with the Outscale BSU builder**.

---

## **:pushpin: How to Test the Examples**

### **Install Packer**
Before running any example, ensure you have **Packer installed**. Follow the official [Packer installation guide](https://developer.hashicorp.com/packer/install) if needed.

### **Set Up Credentials**
You need to configure your Outscale **access keys** and **region** as environment variables:

```bash
export OSC_ACCESS_KEY="myaccesskey"
export OSC_SECRET_KEY="mysecretkey"
```

Alternatively, you can store credentials in a **variables file**.

### **Using a Variables File**
If you prefer, you can define credentials and other parameters in the `variables.auto.pkrvars.hcl` file instead of using environment variables. Simply **edit** the file and add your configuration.

---

## **:pushpin: Running an Example**
Each example is contained in its own folder and can be tested individually.

1. **Initialize the Packer environment** (downloads required plugins):
   ```bash
   packer init .
   ```
   
2. **Validate the configuration** before building:
   ```bash
   packer validate .
   ```
   
3. **Create the Outscale Machine Image (OMI):**
   ```bash
   packer build .
   ```

---

## **:pushpin: Notes**
- Make sure your **credentials are correctly set** before running any example.
- Each folder contains **a separate example**, so navigate to the folder you want to test before running Packer.
- Modify the **`variables.auto.pkrvars.hcl`** file if you need to customize region, instance type, or other parameters.